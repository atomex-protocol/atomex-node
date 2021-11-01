package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
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
	provider   exchange.Exchange
	tracker    *tools.Tracker
	strategies []strategy.Strategy
	symbols    map[string]types.Symbol

	keys              Keys
	atomexMeta        config.Atomex
	quoteProviderMeta QuoteProviderMeta

	orders map[string]Order

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
		atomex.WithWebsocketRestURI(cfg.Atomex.RestAPI))
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

	return &MarketMaker{
		log:               logger.New(logger.WithLogLevel(logLevel), logger.WithModuleName("market_maker")),
		provider:          provider,
		atomex:            atomexExchange,
		tracker:           track,
		strategies:        strategies,
		symbols:           symbols,
		keys:              cfg.Keys,
		atomexMeta:        cfg.Atomex,
		quoteProviderMeta: cfg.QuoteProvider.Meta,
		orders:            make(map[string]Order),
		stop:              make(chan struct{}, 1),
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
func (mm *MarketMaker) Start() error {
	mm.wg.Add(1)
	go mm.listen()

	loadedKeys, err := mm.loadKeys()
	if err != nil {
		return err
	}

	if err := mm.atomex.Connect(loadedKeys); err != nil {
		return errors.Wrap(err, "atomex.Connect")
	}

	providerSymbols := make([]string, 0)
	for symbol := range mm.symbols {
		if s, ok := mm.quoteProviderMeta.FromSymbols[symbol]; ok {
			providerSymbols = append(providerSymbols, s)
		}
	}
	if len(providerSymbols) > 0 {
		if err := mm.provider.Start(providerSymbols...); err != nil {
			return errors.Wrap(err, "quoteProvider.Start")
		}
	}

	if err := mm.tracker.Start(); err != nil {
		return errors.Wrap(err, "tracker.Start")
	}

	if err := mm.sendOneByOneLimits(); err != nil {
		return errors.Wrap(err, "sendOneByOneLimits")
	}

	return nil
}

// Close -
func (mm *MarketMaker) Close() error {
	mm.stop <- struct{}{}
	mm.wg.Wait()

	if err := mm.provider.Close(); err != nil {
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

func (mm *MarketMaker) listen() {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.stop:
			return
		case tick := <-mm.provider.Tickers():
			mm.log.Trace().Str("ask", tick.Ask.String()).Str("bid", tick.Bid.String()).Str("symbol", tick.Symbol).Msg("quote provider's tick")

			if err := mm.sendLimitsByTicker(tick); err != nil {
				mm.log.Err(err).Msg("sendLimitsByTicker")
				continue
			}
		case data := <-mm.atomex.Listen():

			// TODO: check order status and initiate swap

			switch val := data.Value.(type) {
			case atomex.AddOrderResponse:
				mm.log.Info().Int64("order_id", val.OrderID).Msg("atomex new order")
			}

		case err := <-mm.atomex.Errors():
			mm.log.Err(err).Msg("atomex error")

		case operation := <-mm.tracker.Operations():
			mm.log.Info().Str("hash", operation.Hash).Str("chain", operation.ChainType.String()).Str("status", operation.Status.String()).Msg("operation's status changed")

			// TODO: check initiate operation status

		case swap := <-mm.tracker.StatusChanged():
			swap.Log(mm.log.Info()).Msg("swap's status changed")

			// TODO: check swap status and redeem

		}
	}
}

func (mm *MarketMaker) sendOneByOneLimits() error {
	for i := range mm.strategies {
		s, ok := mm.strategies[i].(*strategy.OneByOne)
		if !ok {
			continue
		}

		quotes, err := s.Quotes(nil)
		if err != nil {
			mm.log.Err(err).Msg("get strategy quotes")
			continue
		}

		for i := range quotes {
			symbol, ok := mm.atomexMeta.Symbols[quotes[i].Symbol]
			if !ok {
				continue
			}

			symbolInfo, ok := mm.symbols[quotes[i].Symbol]
			if !ok {
				continue
			}

			receiver, err := mm.getReceiverWallet(symbolInfo, quotes[i].Side)
			if err != nil {
				return errors.Wrap(err, "getReceiverWallet")
			}

			sender, err := mm.getSenderWallet(symbolInfo, quotes[i].Side)
			if err != nil {
				return errors.Wrap(err, "getSenderWallet")
			}

			scrt, err := mm.secret(sender.Private, sender.Address, time.Now().UnixNano())
			if err != nil {
				return errors.Wrap(err, "secret")
			}
			price, _ := quotes[i].Price.Float64()
			qty, _ := quotes[i].Volume.Float64()

			request := atomex.AddOrderRequest{
				ClientOrderID: generateClientOrderID(),
				Symbol:        symbol,
				Price:         price,
				Qty:           qty,
				Side:          mustAtomexSide(quotes[i].Side),
				Type:          atomex.OrderTypeReturn,
				Requisites: &atomex.Requisites{
					BaseCurrencyContract:  symbolInfo.Base.Contract,
					QuoteCurrencyContract: symbolInfo.Quote.Contract,
					SecretHash:            scrt.Hash,
					ReceivingAddress:      receiver.Address,
					RefundAddress:         receiver.Address,
					LockTime:              mm.atomexMeta.Settings.LockTime,
					RewardForRedeem:       mm.atomexMeta.Settings.RewardForRedeem,
				},
			}

			if err := mm.atomex.SendOrder(request); err != nil {
				mm.log.Err(err).Msg("atomex.SendOrder")
				continue
			}

			order := requestToOrder(request, scrt)
			mm.orders[order.ClientID] = order
		}
	}
	return nil
}

func (mm *MarketMaker) sendLimitsByTicker(tick exchange.Ticker) error {
	args := strategy.NewArgs().Ask(tick.Ask).Bid(tick.Bid).AskVolume(tick.AskVolume).BidVolume(tick.BidVolume)
	for i := range mm.strategies {
		quotes, err := mm.strategies[i].Quotes(args)
		if err != nil {
			return errors.Wrap(err, "Quotes")
		}
		for i := range quotes {
			// TODO: check current orders level
			symbol, ok := mm.quoteProviderMeta.ToSymbols[tick.Symbol]
			if !ok {
				return errors.Errorf("unknown provider symbol: %s", tick.Symbol)
			}

			atomexSymbol, ok := mm.atomexMeta.Symbols[symbol]
			if !ok {
				return errors.Errorf("unknown atomex symbol: %s", symbol)
			}

			symbolInfo, ok := mm.symbols[quotes[i].Symbol]
			if !ok {
				continue
			}

			price, _ := quotes[i].Price.Float64()
			qty, _ := quotes[i].Volume.Float64()
			receiver, err := mm.getReceiverWallet(symbolInfo, quotes[i].Side)
			if err != nil {
				return errors.Wrap(err, "getReceiverWallet")
			}
			sender, err := mm.getSenderWallet(symbolInfo, quotes[i].Side)
			if err != nil {
				return errors.Wrap(err, "getSenderWallet")
			}
			scrt, err := mm.secret(sender.Private, sender.Address, time.Now().UnixNano())
			if err != nil {
				return errors.Wrap(err, "secret")
			}

			request := atomex.AddOrderRequest{
				ClientOrderID: generateClientOrderID(),
				Symbol:        atomexSymbol,
				Price:         price,
				Qty:           qty,
				Side:          mustAtomexSide(quotes[i].Side),
				Type:          atomex.OrderTypeReturn,
				Requisites: &atomex.Requisites{
					BaseCurrencyContract:  symbolInfo.Base.Contract,
					QuoteCurrencyContract: symbolInfo.Quote.Contract,
					SecretHash:            scrt.Hash,
					ReceivingAddress:      receiver.Address,
					RefundAddress:         receiver.Address,
					LockTime:              mm.atomexMeta.Settings.LockTime,
					RewardForRedeem:       mm.atomexMeta.Settings.RewardForRedeem,
				},
			}

			if err := mm.atomex.SendOrder(request); err != nil {
				return errors.Wrap(err, "atomex.SendOrder")
			}

			order := requestToOrder(request, scrt)
			mm.orders[order.ClientID] = order
		}
	}
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

type secret struct {
	Value string
	Hash  string
}

func (mm *MarketMaker) secret(key []byte, address string, nonce int64) (secret, error) {
	hash := hmac.New(sha256.New, key)
	nonceBytes := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(nonceBytes, nonce)

	if _, err := hash.Write(append([]byte(address), nonceBytes...)); err != nil {
		return secret{}, err
	}

	scrt := hash.Sum(nil)
	var s secret
	s.Value = hex.EncodeToString(scrt)

	secretHash := sha256.Sum256(scrt)
	s.Hash = hex.EncodeToString(secretHash[:])
	return s, nil
}
