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

	orders     *OrdersMap
	swaps      *SwapsMap
	secrets    *Secrets
	operations map[tools.OperationID]chain.Operation

	activeSwaps []atomex.Swap

	wg sync.WaitGroup
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
			binance.WithRestURL(binance.BaseURLServer2),
			binance.WithWebsocketURL(binance.BaseURLWebsocket),
			binance.WithLogLevel(logLevel),
		)
	default:
		return nil, errors.Errorf("unknown quote provider: %s", cfg.QuoteProvider.Kind)
	}

	atomexExchange, err := atomex.NewExchange(
		atomex.WithLogLevel(logLevel),
		atomex.WithSignature(signers.AlgorithmBlake2bWithEcdsaSecp256k1),
		atomex.WithWebsocketURI(cfg.General.Atomex.WsAPI),
	)
	if err != nil {
		return nil, errors.Wrap(err, "atomex.NewExchange")
	}

	trackerOptions := []tools.TrackerOption{
		tools.WithLogLevel(logLevel),
	}
	if cfg.Restore {
		trackerOptions = append(trackerOptions, tools.WithRestore())
	}
	track, err := tools.NewTracker(cfg.General.Chains, trackerOptions...)
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

		for _, symbol := range cfg.General.Symbols {
			if symbol.Name == s.SymbolName {
				symbols[s.SymbolName] = symbol
				break
			}
		}
	}

	synthetics := make(map[string]synthetic.Synthetic)
	for symbol, cfg := range cfg.QuoteProviderMeta.FromSymbols {
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
			atomex.WithURL(cfg.General.Atomex.RestAPI),
			atomex.WithSignatureAlgorithm(signers.AlgorithmEd25519Blake2b),
		),
		tracker:           track,
		strategies:        strategies,
		symbols:           symbols,
		synthetics:        synthetics,
		keys:              cfg.Keys,
		atomexMeta:        cfg.General.Atomex,
		quoteProviderMeta: cfg.QuoteProviderMeta,
		orders:            NewOrdersMap(),
		swaps:             NewSwapsMap(),
		secrets:           NewSecrets(),
		tickers:           make(map[string]exchange.Ticker),
		operations:        make(map[tools.OperationID]chain.Operation),
		activeSwaps:       make([]atomex.Swap, 0),
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
	go mm.listenAtomex(ctx)

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

	// getting active swaps
	if err := mm.getActiveSwaps(ctx); err != nil {
		return errors.Wrap(err, "getActiveSwaps")
	}

	// init tracker

	mm.wg.Add(1)
	go mm.listenTracker(ctx)

	if err := mm.tracker.Start(ctx); err != nil {
		return errors.Wrap(err, "tracker.Start")
	}

	// init quote provider

	mm.wg.Add(1)
	go mm.listenProvider(ctx)

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

	if err := mm.initialize(ctx); err != nil {
		return errors.Wrap(err, "initialize")
	}

	return nil
}

// Close -
func (mm *MarketMaker) Close(ctx context.Context) error {
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
	return nil
}

func (mm *MarketMaker) getWalletForAsset(asset types.Asset) (chain.Wallet, error) {
	return mm.tracker.Wallet(asset.ChainType())
}

func (mm *MarketMaker) getReceiverWallet(symbol types.Symbol, side strategy.Side) (Wallet, error) {
	switch side {
	case strategy.Ask:
		wallet, err := mm.getWalletForAsset(symbol.Quote)
		if err != nil {
			return Wallet{}, err
		}
		return newWallet(wallet, symbol.Quote), nil
	case strategy.Bid:
		wallet, err := mm.getWalletForAsset(symbol.Base)
		if err != nil {
			return Wallet{}, err
		}
		return newWallet(wallet, symbol.Base), nil
	}
	return Wallet{}, errors.Errorf("unknown side: %v", side)
}

func (mm *MarketMaker) getSenderWallet(symbol types.Symbol, side strategy.Side) (Wallet, error) {
	switch side {
	case strategy.Ask:
		wallet, err := mm.getWalletForAsset(symbol.Base)
		if err != nil {
			return Wallet{}, err
		}
		return newWallet(wallet, symbol.Base), nil
	case strategy.Bid:
		wallet, err := mm.getWalletForAsset(symbol.Base)
		if err != nil {
			return Wallet{}, err
		}
		return newWallet(wallet, symbol.Base), nil
	}
	return Wallet{}, errors.Errorf("unknown side: %v", side)
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
	var offset uint64
	for !end {
		ordersCtx, cancelOrders := context.WithTimeout(ctx, time.Second*5)
		defer cancelOrders()

		orders, err := mm.atomexAPI.Orders(ordersCtx, atomex.OrdersRequest{
			Active: true,
			Sort:   atomex.SortDesc,
			Limit:  limitForAtomexRequest,
			Offset: offset,
		})
		if err != nil {
			return errors.Wrap(err, "atomexAPI.Orders")
		}
		offset += uint64(len(orders))
		end = len(orders) != limitForAtomexRequest

		if err := mm.findDuplicatesOrders(orders); err != nil {
			return errors.Wrap(err, "findDuplicatesOrders")
		}
	}
	return nil
}

func (mm *MarketMaker) initializeSwaps(ctx context.Context) error {

	for i := range mm.activeSwaps {
		if mm.activeSwaps[i].User.Status != atomex.SwapStatusInvolved || mm.activeSwaps[i].CounterParty.Status != atomex.SwapStatusInvolved {
			continue
		}

		swap, err := mm.atomexSwapToInternal(mm.activeSwaps[i])
		if err != nil {
			return err
		}
		internalSwap := mm.swaps.LoadOrStore(swap.HashedSecret, swap)

		if err := mm.restoreSecretForAtomexSwap(ctx, mm.activeSwaps[i], internalSwap); err != nil {
			return errors.Wrap(err, "restoreSecretForAtomexSwap")
		}

		if err := mm.handleAtomexSwapUpdate(mm.activeSwaps[i]); err != nil {
			return errors.Wrap(err, "handleAtomexSwapUpdate")
		}

		if err := mm.initiateInvolvedSwap(ctx, mm.activeSwaps[i]); err != nil {
			return errors.Wrap(err, "initiateInvolvedSwap")
		}
	}

	return nil
}

func (mm *MarketMaker) getActiveSwaps(ctx context.Context) error {
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
		mm.activeSwaps = append(mm.activeSwaps, swaps...)
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

func (mm *MarketMaker) restoreSecretForAtomexSwap(ctx context.Context, atomexSwap atomex.Swap, internalSwap *tools.Swap) error {
	if !atomexSwap.IsInitiator || atomexSwap.Secret != "" {
		return nil
	}

	if len(atomexSwap.User.Trades) == 0 {
		return nil
	}

	orderCtx, cancelCtx := context.WithTimeout(ctx, time.Second*5)
	defer cancelCtx()

	order, err := mm.atomexAPI.Order(orderCtx, atomexSwap.User.Trades[0].OrderID)
	if err != nil {
		return err
	}

	var cid clientOrderID
	if err := cid.parse(order.ClientOrderID); err != nil {
		return err
	}

	var side strategy.Side
	switch atomexSwap.Side {
	case atomex.SideBuy:
		side = strategy.Bid
	case atomex.SideSell:
		side = strategy.Ask
	}

	symbol, ok := mm.atomexMeta.FromSymbols[atomexSwap.Symbol]
	if !ok {
		return errors.Errorf("unknown atomex symbol: %s", atomexSwap.Symbol)
	}

	info, ok := mm.symbols[symbol]
	if !ok {
		return errors.Errorf("unknown symbol: %s", symbol)
	}

	wallet, err := mm.getSenderWallet(info, side)
	if err != nil {
		return err
	}

	scrt, err := mm.secret(wallet.Private, wallet.Address, cid.index)
	if err != nil {
		return err
	}

	mm.secrets.Add(chain.Hex(scrt.Hash), chain.Hex(scrt.Value))
	internalSwap.Secret = chain.Hex(scrt.Value)

	return nil
}
