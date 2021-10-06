package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/ethereum"
	"github.com/atomex-protocol/watch_tower/internal/chain/tezos"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	chainsCount = 2
)

// WatchTower -
type WatchTower struct {
	tezos    *tezos.Tezos
	ethereum *ethereum.Ethereum

	swaps      map[chain.Hex]*Swap
	operations map[OperationID]chain.Operation

	needRedeem bool
	needRefund bool
	retryCount uint

	restoreCounter int32
	stopped        bool
	stop           chan struct{}
	wg             sync.WaitGroup
}

// NewWatchTower -
func NewWatchTower(cfg config.Config) (*WatchTower, error) {
	tezosChain, err := tezos.New(cfg.Tezos)
	if err != nil {
		return nil, err
	}
	eth, err := ethereum.New(cfg.Ethereum)
	if err != nil {
		return nil, err
	}

	wt := &WatchTower{
		tezos:      tezosChain,
		ethereum:   eth,
		retryCount: cfg.RetryCountOnFailedTx,
		swaps:      make(map[chain.Hex]*Swap),
		operations: make(map[OperationID]chain.Operation),
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
	if err := wt.tezos.Init(); err != nil {
		return err
	}

	if err := wt.ethereum.Init(); err != nil {
		return err
	}

	wt.wg.Add(1)
	go wt.listen()

	if restore {
		if err := wt.restore(); err != nil {
			return err
		}
	} else {
		wt.restoreCounter = chainsCount
	}

	if err := wt.tezos.Run(); err != nil {
		return err
	}

	if err := wt.ethereum.Run(); err != nil {
		return err
	}

	return nil
}

func (wt *WatchTower) restore() error {
	if err := wt.tezos.Restore(); err != nil {
		return err
	}

	if err := wt.ethereum.Restore(); err != nil {
		return err
	}

	return nil
}

// Close -
func (wt *WatchTower) Close() error {
	wt.stop <- struct{}{}
	wt.stopped = true
	wt.wg.Wait()

	if err := wt.tezos.Close(); err != nil {
		return err
	}

	if err := wt.ethereum.Close(); err != nil {
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

		// Tezos
		case event := <-wt.tezos.Events():
			if err := wt.onEvent(event); err != nil {
				log.Err(err).Msg("onEvent")
			}
		case operation := <-wt.tezos.Operations():
			if err := wt.processOperation(operation); err != nil {
				log.Err(err).Msg("processOperation")
			}

		// Ethereum
		case event := <-wt.ethereum.Events():
			if err := wt.onEvent(event); err != nil {
				log.Err(err).Msg("onEvent")
			}
		case operation := <-wt.ethereum.Operations():
			if err := wt.processOperation(operation); err != nil {
				log.Err(err).Msg("processOperation")
			}

		// Manager channels
		case <-ticker.C:
			wt.checkRefundTime()
		}
	}
}

func (wt *WatchTower) onEvent(event chain.Event) error {
	switch e := event.(type) {
	case chain.InitEvent:
		if err := wt.onInit(e); err != nil {
			return errors.Wrap(err, "onInit")
		}
	case chain.RedeemEvent:
		if err := wt.onRedeem(e); err != nil {
			return errors.Wrap(err, "onRedeem")
		}
	case chain.RefundEvent:
		if err := wt.onRefund(e); err != nil {
			return errors.Wrap(err, "onRefund")
		}
	case chain.RestoredEvent:
		atomic.AddInt32(&wt.restoreCounter, 1)
		log.Info().Str("blockchain", e.Chain.String()).Msg("restored")

		if wt.restoreCounter == chainsCount {
			for i := range wt.swaps {
				if err := wt.processSwap(wt.swaps[i]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (wt *WatchTower) onInit(event chain.InitEvent) error {
	swap := wt.getSwap(event)
	swap.FromInitEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) onRedeem(event chain.RedeemEvent) error {
	swap := wt.getSwap(event)
	swap.FromRedeemEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) onRefund(event chain.RefundEvent) error {
	swap := wt.getSwap(event)
	swap.FromRefundEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) getSwap(event chain.Event) *Swap {
	s, ok := wt.swaps[event.HashedSecret()]
	if !ok {
		s = NewSwap(event)
		wt.swaps[event.HashedSecret()] = s
	}
	return s
}

func (wt *WatchTower) processSwap(swap *Swap) error {
	swap.log()

	if swap.RetryCount >= wt.retryCount {
		delete(wt.swaps, swap.HashedSecret)
		log.Info().Str("hashed_secret", swap.HashedSecret.String()).Msg("swap retry count transaction exceeded")
		return nil
	}

	switch swap.Status {
	case StatusRedeemedOnce:
		if wt.needRedeem {
			if err := wt.redeem(swap); err != nil {
				return err
			}
		}
	case StatusRefundedOnce:
	case StatusRedeemed, StatusRefunded:
		delete(wt.swaps, swap.HashedSecret)
	default:
	}
	return nil
}

func (wt *WatchTower) checkRefundTime() {
	if !wt.needRefund || wt.restoreCounter < chainsCount {
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
	if wt.restoreCounter < chainsCount {
		return nil
	}

	if leg := swap.Leg(); leg != nil {
		swap.RetryCount++
		switch leg.ChainType {
		case chain.ChainTypeEthereum:
			if err := wt.ethereum.Redeem(swap.HashedSecret, swap.Secret, leg.Contract); err != nil {
				return err
			}
			time.Sleep(time.Second)
		case chain.ChainTypeTezos:
			if err := wt.tezos.Redeem(swap.HashedSecret, swap.Secret, leg.Contract); err != nil {
				return err
			}
			time.Sleep(time.Second)
		default:
			return errors.Errorf("unknown chain type: %v", leg.ChainType)
		}
	}

	return nil
}

func (wt *WatchTower) refund(swap *Swap) error {
	if wt.restoreCounter < chainsCount {
		return nil
	}

	if leg := swap.Leg(); leg != nil {
		swap.RetryCount++
		switch leg.ChainType {
		case chain.ChainTypeEthereum:
			if err := wt.ethereum.Refund(swap.HashedSecret, leg.Contract); err != nil {
				return err
			}
			time.Sleep(time.Second)
		case chain.ChainTypeTezos:
			if err := wt.tezos.Refund(swap.HashedSecret, leg.Contract); err != nil {
				return err
			}
			time.Sleep(time.Second)
		default:
			return errors.Errorf("unknown chain type: %v", leg.ChainType)
		}
	}

	return nil
}

func (wt *WatchTower) processOperation(operation chain.Operation) error {
	id := OperationID{operation.Hash, operation.ChainType}

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
				return wt.processSwap(swap)
			}
		}
	}

	return nil
}
