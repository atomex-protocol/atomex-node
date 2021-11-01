package main

import (
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

	stopped bool
	stop    chan struct{}
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
	track, err := tools.NewTracker(cfg.Chains, opts...)
	if err != nil {
		return nil, err
	}
	wt := &WatchTower{
		tracker:    track,
		retryCount: cfg.RetryCountOnFailedTx,
		operations: make(map[tools.OperationID]chain.Operation),
		swaps:      make(map[chain.Hex]*Swap),
		stop:       make(chan struct{}, 1),
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

// Run -
func (wt *WatchTower) Run(restore bool) error {
	wt.wg.Add(1)
	go wt.listen()

	if err := wt.tracker.Start(); err != nil {
		return err
	}

	return nil
}

// Close -
func (wt *WatchTower) Close() error {
	wt.stop <- struct{}{}
	wt.stopped = true
	wt.wg.Wait()

	if err := wt.tracker.Close(); err != nil {
		return err
	}

	close(wt.stop)
	return nil
}

func (wt *WatchTower) listen() {
	defer wt.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-wt.stop:
			return

		// Tracker
		case swap := <-wt.tracker.StatusChanged():
			swap.Log(log.Info()).Msg("swap info")

			s, ok := wt.swaps[swap.HashedSecret]
			if !ok {
				s = &Swap{swap, 0}
				wt.swaps[swap.HashedSecret] = s
			}

			if err := wt.onSwap(s); err != nil {
				log.Err(err).Msg("onSwap")
			}
		case operation := <-wt.tracker.Operations():
			if err := wt.onOperation(operation); err != nil {
				log.Err(err).Msg("onOperation")
			}

		// Manager channels
		case <-ticker.C:
			wt.checkRefundTime()
		}
	}
}

func (wt *WatchTower) onSwap(swap *Swap) error {
	if swap.RetryCount >= wt.retryCount {
		delete(wt.swaps, swap.HashedSecret)
		log.Info().Str("hashed_secret", swap.HashedSecret.String()).Msg("swap retry count transaction exceeded")
		return nil
	}

	switch swap.Status {
	case tools.StatusRedeemedOnce:
		if wt.needRedeem {
			if err := wt.redeem(swap); err != nil {
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

func (wt *WatchTower) checkRefundTime() {
	if !wt.needRefund {
		return
	}

	for hashedSecret, swap := range wt.swaps {
		if wt.stopped {
			return
		}
		if swap.RefundTime.UTC().Before(time.Now().UTC()) {
			if err := wt.refund(swap); err != nil {
				log.Err(err).Msg("refund")
				continue
			}
			delete(wt.swaps, hashedSecret)
		}
	}
}

func (wt *WatchTower) redeem(swap *Swap) error {
	if leg := swap.Leg(); leg != nil {
		swap.RetryCount++
		return wt.tracker.Redeem(swap.Swap, *leg)
	}

	return nil
}

func (wt *WatchTower) refund(swap *Swap) error {
	if leg := swap.Leg(); leg != nil {
		swap.RetryCount++
		return wt.tracker.Refund(swap.Swap, *leg)
	}

	return nil
}

func (wt *WatchTower) onOperation(operation chain.Operation) error {
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
				return wt.onSwap(swap)
			}
		}
	}

	return nil
}
