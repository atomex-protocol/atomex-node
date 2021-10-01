package main

import (
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/ethereum"
	"github.com/atomex-protocol/watch_tower/internal/chain/tezos"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

	restoring bool
	stopped   bool
	stop      chan struct{}
	wg        sync.WaitGroup
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
	log.Info().Str("blockchain", "tezos").Msg("running...")
	if err := wt.tezos.Run(); err != nil {
		return err
	}

	log.Info().Str("blockchain", "ethereum").Msg("running...")
	if err := wt.ethereum.Run(); err != nil {
		return err
	}

	wt.wg.Add(1)
	go wt.listen()

	if restore {
		if err := wt.restore(); err != nil {
			return err
		}
	}

	return nil
}

func (wt *WatchTower) restore() error {
	defer func() {
		wt.restoring = false
	}()
	wt.restoring = true

	if err := wt.tezos.Restore(); err != nil {
		return err
	}

	if err := wt.ethereum.Restore(); err != nil {
		return err
	}

	time.Sleep(time.Second * 30)

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

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-wt.stop:
			return

		// Tezos
		case event := <-wt.tezos.InitEvents():
			if err := wt.onInit(event); err != nil {
				log.Err(err).Msg("onInit")
			}
		case event := <-wt.tezos.RedeemEvents():
			if err := wt.onRedeem(event); err != nil {
				log.Err(err).Msg("onRedeem")
			}
		case event := <-wt.tezos.RefundEvents():
			if err := wt.onRefund(event); err != nil {
				log.Err(err).Msg("onRefund")
			}
		case operation := <-wt.tezos.Operations():
			if err := wt.processOperation(operation); err != nil {
				log.Err(err).Msg("processOperation")
			}

		// Ethereum
		case event := <-wt.ethereum.InitEvents():
			if err := wt.onInit(event); err != nil {
				log.Err(err).Msg("onInit")
			}
		case event := <-wt.ethereum.RedeemEvents():
			if err := wt.onRedeem(event); err != nil {
				log.Err(err).Msg("onRedeem")
			}
		case event := <-wt.ethereum.RefundEvents():
			if err := wt.onRefund(event); err != nil {
				log.Err(err).Msg("onRefund")
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

func (wt *WatchTower) onInit(event chain.InitEvent) error {
	swap := wt.getSwap(event.Event)
	swap.FromInitEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) onRedeem(event chain.RedeemEvent) error {
	swap := wt.getSwap(event.Event)
	swap.FromRedeemEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) onRefund(event chain.RefundEvent) error {
	swap := wt.getSwap(event.Event)
	swap.FromRefundEvent(event)
	return wt.processSwap(swap)
}

func (wt *WatchTower) getSwap(event chain.Event) *Swap {
	s, ok := wt.swaps[event.HashedSecret]
	if !ok {
		s = NewSwap(event)
		wt.swaps[event.HashedSecret] = s
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
	if !wt.needRefund || wt.restoring {
		return
	}

	for hashedSecret, swap := range wt.swaps {
		if wt.stopped {
			return
		}
		if swap.RefundTime.UTC().Before(time.Now().UTC()) {
			if err := wt.refund(swap); err != nil {
				log.Err(err).Msg("checkRefundTime")
				continue
			}
			delete(wt.swaps, hashedSecret)
			time.Sleep(time.Second)
		}
	}
}

func (wt *WatchTower) redeem(swap *Swap) error {
	if swap.Acceptor.ChainType == chain.ChainTypeUnknown {
		swap.RetryCount = wt.retryCount
		delete(wt.swaps, swap.HashedSecret)
		return nil
	}

	if wt.restoring {
		return nil
	}

	log.Info().Str("hashed_secret", swap.HashedSecret.String()).Str("blockchain", swap.Initiator.ChainType.String()).Msg("redeem")

	swap.RetryCount++
	switch swap.Initiator.ChainType {
	case chain.ChainTypeEthereum:
		return wt.ethereum.Redeem(swap.HashedSecret, swap.Secret, swap.Contract)
	case chain.ChainTypeTezos:
		return wt.tezos.Redeem(swap.HashedSecret, swap.Secret, swap.Contract)
	default:
		return errors.Errorf("unknown chain type: %v", swap.Initiator.ChainType)
	}
}

func (wt *WatchTower) refund(swap *Swap) error {
	if swap.Acceptor.ChainType == chain.ChainTypeUnknown {
		swap.RetryCount = wt.retryCount
		delete(wt.swaps, swap.HashedSecret)
		return nil
	}

	if wt.restoring {
		return nil
	}

	log.Info().Str("hashed_secret", swap.HashedSecret.String()).Str("blockchain", swap.Initiator.ChainType.String()).Msg("refund")

	swap.RetryCount++
	switch swap.Initiator.ChainType {
	case chain.ChainTypeEthereum:
		return wt.ethereum.Refund(swap.HashedSecret, swap.Contract)
	case chain.ChainTypeTezos:
		return wt.tezos.Refund(swap.HashedSecret, swap.Contract)
	default:
		return errors.Errorf("unknown chain type: %v", swap.Initiator.ChainType)
	}
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
