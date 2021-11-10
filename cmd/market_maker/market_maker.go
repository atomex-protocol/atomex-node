package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/synthetic"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/atomex-protocol/watch_tower/internal/exchange/binance"
	"github.com/atomex-protocol/watch_tower/internal/keys"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// MarketMaker -
type MarketMaker struct {
	log zerolog.Logger

	atomex     *atomex.Exchange
	atomexAPI  *atomex.Rest
	provider   exchange.Exchange
	tracker    *tools.Tracker
	strategies []strategy.Strategy
	symbols    map[string]types.Symbol

	keys              Keys
	atomexMeta        config.Atomex
	quoteProviderMeta QuoteProviderMeta
	tickers           map[string]exchange.Ticker
	synthetics        map[string]synthetic.Synthetic

	ordersCounter uint64

	orders     *OrdersMap
	swaps      *SwapsMap
	operations map[tools.OperationID]chain.Operation

	wg   sync.WaitGroup
	stop chan struct{}
}

// NewMarketMaker -
func NewMarketMaker(cfg Config) (*MarketMaker, error) {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "zerolog.ParseLevel")
	}

	var provider exchange.Exchange
	switch cfg.QuoteProvider.Kind {
	case QuoteProviderKindBinance:
		provider = binance.NewBinance(
			binance.WithRestURL(binance.BaseURLServer1),
			binance.WithWebsocketURL(binance.BaseURLWebsocket),
			binance.WithLogLevel(logLevel),
		)
	default:
		return nil, errors.Errorf("unknown quote provider: %s", cfg.QuoteProvider.Kind)
	}

	atomexExchange, err := atomex.NewExchange(
		atomex.WithLogLevel(logLevel),
		atomex.WithSignature(signers.AlgorithmBlake2bWithEcdsaSecp256k1),
		atomex.WithWebsocketURI(cfg.Atomex.WsAPI),
	)
	if err != nil {
		return nil, errors.Wrap(err, "atomex.NewExchange")
	}

	track, err := tools.NewTracker(cfg.Chains, tools.WithLogLevel(logLevel)) //, tools.WithRestore())
	if err != nil {
		return nil, err
	}

	symbols := make(map[string]types.Symbol)
	strategies := make([]strategy.Strategy, 0)
	for _, s := range cfg.Strategies {
		strategy, err := strategy.New(s)
		if err != nil {
			return nil, err
		}
		strategies = append(strategies, strategy)

		for _, symbol := range cfg.Info.Symbols {
			if symbol.Name == s.SymbolName {
				symbols[s.SymbolName] = symbol
				break
			}
		}
	}

	synthetics := make(map[string]synthetic.Synthetic)
	for symbol, cfg := range cfg.QuoteProvider.Meta.FromSymbols {
		synth, err := synthetic.New(symbol, cfg)
		if err != nil {
			return nil, err
		}
		synthetics[symbol] = synth
	}

	return &MarketMaker{
		log:      logger.New(logger.WithLogLevel(logLevel), logger.WithModuleName("market_maker")),
		provider: provider,
		atomex:   atomexExchange,
		atomexAPI: atomex.NewRest(
			atomex.WithURL(cfg.Atomex.RestAPI),
			atomex.WithSignatureAlgorithm(signers.AlgorithmBlake2bWithEcdsaSecp256k1),
		),
		tracker:           track,
		strategies:        strategies,
		symbols:           symbols,
		synthetics:        synthetics,
		keys:              cfg.Keys,
		atomexMeta:        cfg.Atomex,
		quoteProviderMeta: cfg.QuoteProvider.Meta,
		orders:            NewOrdersMap(),
		swaps:             NewSwapsMap(),
		tickers:           make(map[string]exchange.Ticker),
		operations:        make(map[tools.OperationID]chain.Operation),
		stop:              make(chan struct{}, 3),
	}, nil
}

func (mm *MarketMaker) loadKeys() (*signers.Key, error) {
	keysStorage, err := keys.New(mm.keys.Kind)
	if err != nil {
		return nil, err
	}

	loadedKeys, err := keysStorage.Get(mm.keys.File)
	if err == nil {
		return loadedKeys, err
	}

	if os.IsNotExist(err) && mm.keys.GenerateIfNotExists {
		return keysStorage.Create(mm.keys.File)
	}

	return nil, err
}

// Start -
func (mm *MarketMaker) Start(ctx context.Context) error {
	mm.wg.Add(1)
	go mm.listenAtomex()

	loadedKeys, err := mm.loadKeys()
	if err != nil {
		return err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := mm.atomexAPI.Auth(ctxTimeout, loadedKeys); err != nil {
		return errors.Wrap(err, "Token")
	}

	if err := mm.atomex.Connect(atomex.TokenResponse{
		Token: mm.atomexAPI.GetToken(),
	}); err != nil {
		return errors.Wrap(err, "atomex.Connect")
	}

	if err := mm.initialize(ctx); err != nil {
		return errors.Wrap(err, "initialize")
	}

	// init tracker

	mm.wg.Add(1)
	go mm.listenTracker()

	if err := mm.tracker.Start(); err != nil {
		return errors.Wrap(err, "tracker.Start")
	}

	// init quote provider

	mm.wg.Add(1)
	go mm.listenProvider()

	providerSymbols := make([]string, 0)
	for symbol := range mm.symbols {
		if s, ok := mm.quoteProviderMeta.FromSymbols[symbol]; ok {
			providerSymbols = append(providerSymbols, s.Symbols...)
		}
	}
	if len(providerSymbols) > 0 {
		if err := mm.provider.Start(providerSymbols...); err != nil {
			return errors.Wrap(err, "quoteProvider.Start")
		}
	}

	return nil
}

// Close -
func (mm *MarketMaker) Close(ctx context.Context) error {
	for i := 0; i < cap(mm.stop); i++ {
		mm.stop <- struct{}{}
	}
	mm.wg.Wait()

	if err := mm.provider.Close(); err != nil {
		return err
	}

	if err := mm.cancelAll(ctx); err != nil {
		return err
	}
	if err := mm.atomex.Close(); err != nil {
		return err
	}
	if err := mm.tracker.Close(); err != nil {
		return err
	}

	close(mm.stop)
	return nil
}

func (mm *MarketMaker) getWalletForAsset(asset types.Asset) (chain.Wallet, error) {
	return mm.tracker.Wallet(asset.ChainType())
}

func (mm *MarketMaker) getReceiverWallet(symbol types.Symbol, side strategy.Side) (chain.Wallet, error) {
	switch side {
	case strategy.Ask:
		return mm.getWalletForAsset(symbol.Quote)
	case strategy.Bid:
		return mm.getWalletForAsset(symbol.Base)
	}
	return chain.Wallet{}, errors.Errorf("unknown side: %v", side)
}

func (mm *MarketMaker) getSenderWallet(symbol types.Symbol, side strategy.Side) (chain.Wallet, error) {
	switch side {
	case strategy.Ask:
		return mm.getWalletForAsset(symbol.Base)
	case strategy.Bid:
		return mm.getWalletForAsset(symbol.Quote)
	}
	return chain.Wallet{}, errors.Errorf("unknown side: %v", side)
}

const limitForAtomexRequest = 100

func (mm *MarketMaker) initialize(ctx context.Context) error {
	if err := mm.initializeOrders(ctx); err != nil {
		return errors.Wrap(err, "initializeOrders")
	}

	if err := mm.initializeSwaps(ctx); err != nil {
		return errors.Wrap(err, "initializeSwaps")
	}

	if err := mm.sendOneByOneLimits(); err != nil {
		return errors.Wrap(err, "sendOneByOneLimits")
	}

	return nil
}

func (mm *MarketMaker) initializeOrders(ctx context.Context) error {
	var end bool
	for !end {
		ordersCtx, cancelOrders := context.WithTimeout(ctx, time.Second*5)
		defer cancelOrders()

		orders, err := mm.atomexAPI.Orders(ordersCtx, atomex.OrdersRequest{
			Active: true,
			Sort:   atomex.SortDesc,
			Limit:  limitForAtomexRequest,
		})
		if err != nil {
			return errors.Wrap(err, "atomexAPI.Orders")
		}
		end = len(orders) != limitForAtomexRequest

		if err := mm.findDuplicatesOrders(orders); err != nil {
			return errors.Wrap(err, "findDuplicatesOrders")
		}
	}
	return nil
}

func (mm *MarketMaker) initializeSwaps(ctx context.Context) error {
	var end bool
	for !end {
		swapsCtx, cancelSwaps := context.WithTimeout(ctx, time.Second*5)
		defer cancelSwaps()

		swaps, err := mm.atomexAPI.Swaps(swapsCtx, atomex.SwapsRequest{
			Active: true,
			Sort:   atomex.SortDesc,
			Limit:  limitForAtomexRequest,
		})
		if err != nil {
			return errors.Wrap(err, "atomexAPI.Swaps")
		}
		end = len(swaps) != limitForAtomexRequest

		for i := range swaps {
			swap, err := mm.atomexSwapToInternal(swaps[i])
			if err != nil {
				return errors.Wrap(err, "atomexSwapToInternal")
			}
			if swap != nil {
				mm.swaps.Store(chain.Hex(swaps[i].SecretHash), swap)
			}
		}
	}
	return nil
}

func (mm *MarketMaker) atomexSwapToInternal(swap atomex.Swap) (*tools.Swap, error) {
	symbol, ok := mm.atomexMeta.FromSymbols[swap.Symbol]
	if !ok {
		return nil, errors.Errorf("unknown atomex symbol: %s", swap.Symbol)
	}

	info, ok := mm.symbols[symbol]
	if !ok {
		return nil, errors.Errorf("unknown symbol: %s", symbol)
	}

	var initiator, acceptor types.Asset
	switch swap.Side {
	case atomex.SideBuy:
		initiator = info.Quote
		acceptor = info.Base
	case atomex.SideSell:
		initiator = info.Base
		acceptor = info.Quote
	}

	initiatorStatus := atomexStatusToInternal(swap.User.Status)
	acceptorStatus := atomexStatusToInternal(swap.CounterParty.Status)

	status := tools.StatusEmpty

	switch {
	case initiatorStatus == tools.StatusRedeemed && acceptorStatus == tools.StatusRedeemed:
		status = tools.StatusRedeemed
	case initiatorStatus == tools.StatusRedeemed || acceptorStatus == tools.StatusRedeemed:
		status = tools.StatusRedeemedOnce
	case initiatorStatus == tools.StatusRefunded && acceptorStatus == tools.StatusRefunded:
		status = tools.StatusRefunded
	case initiatorStatus == tools.StatusRefunded || acceptorStatus == tools.StatusRefunded:
		status = tools.StatusRefundedOnce
	case initiatorStatus == tools.StatusInitiated && acceptorStatus == tools.StatusInitiated:
		status = tools.StatusInitiated
	case initiatorStatus == tools.StatusInitiated || acceptorStatus == tools.StatusInitiated:
		status = tools.StatusInitiatedOnce
	}

	initiatorWallet, err := mm.getWalletForAsset(initiator)
	if err != nil {
		return nil, err
	}
	acceptorWallet, err := mm.getWalletForAsset(acceptor)
	if err != nil {
		return nil, err
	}

	return &tools.Swap{
		HashedSecret: chain.Hex(swap.SecretHash),
		Secret:       chain.Hex(swap.Secret),
		Status:       status,
		RefundTime:   swap.TimeStamp.Add(time.Duration(swap.User.Requisites.LockTime) * time.Second),
		Symbol:       info,
		Initiator: tools.Leg{
			Status:    initiatorStatus,
			ChainType: initiator.ChainType(),
			Contract:  initiator.Contract,
			Address:   initiatorWallet.Address,
		},
		Acceptor: tools.Leg{
			Status:    acceptorStatus,
			ChainType: acceptor.ChainType(),
			Contract:  acceptor.Contract,
			Address:   acceptorWallet.Address,
		},
	}, nil
}
