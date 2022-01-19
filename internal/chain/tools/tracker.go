package tools

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/ethereum"
	"github.com/atomex-protocol/watch_tower/internal/chain/tezos"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// tracker -
type Tracker struct {
	tezos *tezos.Tezos
	eth   *ethereum.Ethereum

	logger zerolog.Logger

	restoreCounter int32
	needRestore    bool

	swaps         map[chain.Hex]*Swap
	statusChanged chan Swap
	operations    chan chain.Operation

	wg sync.WaitGroup
}

// NewTracker -
func NewTracker(cfg Config, opts ...TrackerOption) (*Tracker, error) {
	tezosChain, err := tezos.New(tezos.Config{
		Node:            cfg.Tezos.Node,
		TzKT:            cfg.Tezos.TzKT,
		Contract:        cfg.Tezos.Contract,
		Tokens:          cfg.Tezos.Tokens,
		MinPayOff:       cfg.Tezos.MinPayOff,
		TTL:             cfg.Tezos.TTL,
		OperaitonParams: cfg.Tezos.OperaitonParams,
		LogLevel:        zerolog.InfoLevel,
	})
	if err != nil {
		return nil, err
	}

	eth, err := ethereum.New(ethereum.Config{
		EthContract:   cfg.Ethereum.EthAddress,
		Erc20Contract: cfg.Ethereum.Erc20Address,
		NodeURL:       cfg.Ethereum.Node,
		WssURL:        cfg.Ethereum.Wss,
		MinPayOff:     cfg.Ethereum.MinPayOff,
		LogLevel:      zerolog.InfoLevel,
	})
	if err != nil {
		return nil, err
	}

	t := &Tracker{
		tezos: tezosChain,
		eth:   eth,

		logger: logger.New(logger.WithModuleName("tracker")),

		swaps:         make(map[chain.Hex]*Swap),
		operations:    make(chan chain.Operation, 1024),
		statusChanged: make(chan Swap, 1024),
	}
	for i := range opts {
		opts[i](t)
	}
	return t, nil
}

// StatusChanged -
func (t *Tracker) StatusChanged() <-chan Swap {
	return t.statusChanged
}

// Operations -
func (t *Tracker) Operations() <-chan chain.Operation {
	return t.operations
}

// Close -
func (t *Tracker) Close() error {
	t.wg.Wait()

	if err := t.eth.Close(); err != nil {
		return err
	}
	if err := t.tezos.Close(); err != nil {
		return err
	}

	close(t.operations)
	close(t.statusChanged)
	return nil
}

// Start -
func (t *Tracker) Start(ctx context.Context) error {
	if err := t.tezos.Init(ctx); err != nil {
		return err
	}

	if err := t.eth.Init(ctx); err != nil {
		return err
	}

	t.wg.Add(1)
	go t.listen(ctx)

	if t.needRestore {
		if err := t.restore(ctx); err != nil {
			return err
		}
	} else {
		t.restoreCounter = chainsCount
	}

	if err := t.tezos.Run(ctx); err != nil {
		return err
	}

	if err := t.eth.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (t *Tracker) restore(ctx context.Context) error {
	if err := t.tezos.Restore(ctx); err != nil {
		return err
	}

	if err := t.eth.Restore(); err != nil {
		return err
	}

	return nil
}

func (t *Tracker) listen(ctx context.Context) {
	defer t.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		// Tezos
		case event := <-t.tezos.Events():
			t.onEvent(event)
		case operation := <-t.tezos.Operations():
			t.operations <- operation

		// Ethereum
		case event := <-t.eth.Events():
			t.onEvent(event)
		case operation := <-t.eth.Operations():
			t.operations <- operation
		}
	}
}

func (t *Tracker) onEvent(event chain.Event) {
	switch e := event.(type) {
	case chain.InitEvent:
		t.onInit(e)
	case chain.RedeemEvent:
		t.onRedeem(e)
	case chain.RefundEvent:
		t.onRefund(e)
	case chain.RestoredEvent:
		t.logger.Info().Str("blockchain", e.Chain.String()).Msg("restored")
		atomic.AddInt32(&t.restoreCounter, 1)

		if t.restoreCounter == chainsCount {
			for id := range t.swaps {
				t.statusChanged <- *t.swaps[id]
			}
		}
	}
}

func (t *Tracker) onInit(event chain.InitEvent) {
	swap := t.getSwap(event)
	swap.fromInitEvent(event)
	if t.restoreCounter == chainsCount {
		t.statusChanged <- *swap
	}
}

func (t *Tracker) onRedeem(event chain.RedeemEvent) {
	swap := t.getSwap(event)
	swap.fromRedeemEvent(event)
	if t.restoreCounter == chainsCount {
		t.statusChanged <- *swap
	}
}

func (t *Tracker) onRefund(event chain.RefundEvent) {
	swap := t.getSwap(event)
	swap.fromRefundEvent(event)
	if t.restoreCounter == chainsCount {
		t.statusChanged <- *swap
	}
}

func (t *Tracker) getSwap(event chain.Event) *Swap {
	s, ok := t.swaps[event.HashedSecret()]
	if !ok {
		s = NewSwap(event)
		t.swaps[event.HashedSecret()] = s
	}
	return s
}

// Redeem -
func (t *Tracker) Redeem(ctx context.Context, swap Swap, leg Leg) error {
	switch leg.ChainType {
	case chain.ChainTypeEthereum:
		if err := t.eth.Redeem(ctx, swap.HashedSecret, swap.Secret, leg.Contract); err != nil {
			return err
		}
	case chain.ChainTypeTezos:
		if err := t.tezos.Redeem(ctx, swap.HashedSecret, swap.Secret, leg.Contract); err != nil {
			return err
		}
	default:
		return errors.Wrapf(ErrUnknownChainType, "Redeem %v", leg.ChainType)
	}
	return nil
}

// Redeem -
func (t *Tracker) Refund(ctx context.Context, swap Swap, leg Leg) error {
	switch leg.ChainType {
	case chain.ChainTypeEthereum:
		if err := t.eth.Refund(ctx, swap.HashedSecret, leg.Contract); err != nil {
			return err
		}
	case chain.ChainTypeTezos:
		if err := t.tezos.Refund(ctx, swap.HashedSecret, leg.Contract); err != nil {
			return err
		}
	default:
		return errors.Wrapf(ErrUnknownChainType, "Refund %v", leg.ChainType)
	}
	return nil
}

// Initiate -
func (t *Tracker) Initiate(ctx context.Context, args chain.InitiateArgs, chainType chain.ChainType) error {
	switch chainType {
	case chain.ChainTypeEthereum:
		if err := t.eth.Initiate(ctx, args); err != nil {
			return err
		}
	case chain.ChainTypeTezos:
		if err := t.tezos.Initiate(ctx, args); err != nil {
			return err
		}
	default:
		return errors.Wrapf(ErrUnknownChainType, "Initiate %v", chainType)
	}
	return nil
}

// Wallet -
func (t *Tracker) Wallet(typ chain.ChainType) (chain.Wallet, error) {
	switch typ {
	case chain.ChainTypeEthereum:
		return t.eth.Wallet(), nil
	case chain.ChainTypeTezos:
		return t.tezos.Wallet(), nil
	}
	return chain.Wallet{}, errors.Wrapf(ErrUnknownChainType, typ.String())
}
