package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// WatchTower -
type WatchTower struct {
	tracker    *tools.Tracker
	operations map[tools.OperationID]chain.Operation
	swaps      map[chain.Hex]*Swap

	needRedeem bool
	needRefund bool
	retryCount uint
	uptimeAPI  string

	stopped bool
	wg      sync.WaitGroup
}

// NewWatchTower -
func NewWatchTower(cfg Config) (*WatchTower, error) {
	opts := []tools.TrackerOption{
		tools.WithLogLevel(zerolog.InfoLevel),
	}
	if cfg.Restore {
		opts = append(opts, tools.WithRestore())
	}
	track, err := tools.NewTracker(cfg.General.Chains, opts...)
	if err != nil {
		return nil, err
	}
	wt := &WatchTower{
		tracker:    track,
		retryCount: cfg.RetryCountOnFailedTx,
		uptimeAPI:  cfg.General.Atomex.UptimeAPI,
		operations: make(map[tools.OperationID]chain.Operation),
		swaps:      make(map[chain.Hex]*Swap),
	}
	if wt.retryCount == 0 {
		wt.retryCount = 3
	}

	for i := range cfg.Types {
		switch cfg.Types[i] {
		case "redeem":
			wt.needRedeem = true
		case "refund":
			wt.needRefund = true
		}
	}

	return wt, nil
}

const minus30Minutes = -30 * time.Minute

// Run -
func (wt *WatchTower) Run(ctx context.Context, restore bool) error {
	wt.wg.Add(1)
	go wt.listen(ctx)

	if err := wt.tracker.Start(ctx); err != nil {
		return err
	}

	return nil
}

// Close -
func (wt *WatchTower) Close() error {
	wt.stopped = true
	wt.wg.Wait()

	if err := wt.tracker.Close(); err != nil {
		return err
	}

	return nil
}

func (wt *WatchTower) listen(ctx context.Context) {
	defer wt.wg.Done()

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	heartbeatTicker := time.NewTicker(time.Hour)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		// Tracker
		case swap := <-wt.tracker.StatusChanged():
			swap.Log(log.Info()).Msg("swap info")

			s, ok := wt.swaps[swap.HashedSecret]
			if !ok {
				s = &Swap{swap, 0}
				wt.swaps[swap.HashedSecret] = s
			} else {
				s.merge(swap)
			}

			if err := wt.onSwap(ctx, s); err != nil {
				log.Err(err).Msg("onSwap")
			}
		case operation := <-wt.tracker.Operations():
			if err := wt.onOperation(ctx, operation); err != nil {
				log.Err(err).Msg("onOperation")
			}

		// Manager channels
		case <-ticker.C:
			wt.checkNextActionTime(ctx)

		case <-heartbeatTicker.C:
			wt.heartbeat()
		}
	}
}

func (wt *WatchTower) onSwap(ctx context.Context, swap *Swap) error {
	if swap.RetryCount >= wt.retryCount {
		delete(wt.swaps, swap.HashedSecret)
		log.Info().Str("hashed_secret", swap.HashedSecret.String()).Msg("swap retry count transaction exceeded")
		return nil
	}

	switch swap.Status {
	case tools.StatusRedeemedOnce:
		if wt.needRedeem {
			if err := wt.redeem(ctx, swap); err != nil {
				return err
			}
		}
	case tools.StatusRefundedOnce:
	case tools.StatusRedeemed, tools.StatusRefunded:
		delete(wt.swaps, swap.HashedSecret)
	default:
	}
	return nil
}

func (wt *WatchTower) checkNextActionTime(ctx context.Context) {
	if !wt.needRefund && !wt.needRedeem {
		return
	}

	for hashedSecret, swap := range wt.swaps {
		if wt.stopped {
			return
		}

		if wt.needRedeem && swap.Status == tools.StatusRedeemedOnce {
			if err := wt.redeem(ctx, swap); err != nil {
				log.Err(err).Msg("redeem")
			}
		}

		if wt.needRefund {
			if swap.IsUnknown() {
				continue
			}

			if swap.RefundTime.UTC().Before(time.Now().UTC()) {
				if err := wt.refund(ctx, swap); err != nil {
					log.Err(err).Msg("refund")
					continue
				}

				delete(wt.swaps, hashedSecret)
			}
		}
	}
}

func (wt *WatchTower) redeem(ctx context.Context, swap *Swap) error {
	utcNow := time.Now().UTC()

	if leg := swap.Leg(); leg != nil && utcNow.Before(swap.RefundTime.UTC()) {
		if swap.RewardForRedeem.IsPositive() {
			swap.RetryCount++
			return wt.tracker.Redeem(ctx, swap.Swap, *leg)
		}

		if swap.RewardForRedeem.IsZero() && utcNow.After(swap.RefundTime.Add(minus30Minutes).UTC()) {
			log.Info().Msg("WatchTower starts redeem for swap with zero reward")
			swap.RetryCount++
			return wt.tracker.Redeem(ctx, swap.Swap, *leg)
		}
	}

	return nil
}

func (wt *WatchTower) refund(ctx context.Context, swap *Swap) error {
	if leg := swap.Leg(); leg != nil {
		swap.RetryCount++
		return wt.tracker.Refund(ctx, swap.Swap, *leg)
	}

	if swap.Acceptor.Status == tools.StatusInitiated && swap.Initiator.Status == tools.StatusInitiated {
		swap.RetryCount++
		if err := wt.tracker.Refund(ctx, swap.Swap, swap.Initiator); err != nil {
			return err
		}
		return wt.tracker.Refund(ctx, swap.Swap, swap.Acceptor)
	}

	return nil
}

func (wt *WatchTower) onOperation(ctx context.Context, operation chain.Operation) error {
	id := tools.OperationID{
		Hash:  operation.Hash,
		Chain: operation.ChainType,
	}

	switch operation.Status {
	case chain.Pending:
		wt.operations[id] = operation
		log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", operation.HashedSecret.String()).Msg("transaction")
	case chain.Applied:
		if old, ok := wt.operations[id]; ok {
			log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", old.HashedSecret.String()).Msg("transaction")
			delete(wt.operations, id)
		}
	case chain.Failed:
		if old, ok := wt.operations[id]; ok {
			log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", old.HashedSecret.String()).Msg("transaction")
			delete(wt.operations, id)

			if swap, ok := wt.swaps[old.HashedSecret]; ok {
				return wt.onSwap(ctx, swap)
			}
		}
	}

	return nil
}

func (wt *WatchTower) heartbeat() {
	requestUri := fmt.Sprintf("%s?msg=OK,%%20%d%%20swaps", wt.uptimeAPI, len(wt.swaps))
	res, err := http.Head(requestUri)

	if err != nil {
		log.Err(err).Msgf("WatchTower 'stay alive' heartbeat failed to be sent")
	} else if res != nil && res.StatusCode != http.StatusOK {
		var location string

		if url, errLoc := res.Location(); errLoc != nil {
			location = url.String()
		}

		log.Err(err).Msgf("WatchTower 'stay alive' heartbeat failed to be sent; response: { code: %v, location: %v }", res.StatusCode, location)
	}
}
