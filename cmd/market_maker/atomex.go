package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/synthetic"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (mm *MarketMaker) listenAtomex(ctx context.Context) {
	defer mm.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case data := <-mm.atomex.Listen():
			if err := mm.handleAtomexUpdate(ctx, data); err != nil {
				mm.log.Err(err).Msg("handleAtomexUpdate")
				continue
			}

		case err := <-mm.atomex.Errors():
			mm.log.Err(err).Msg("atomex error")
		}
	}
}

func (mm *MarketMaker) handleAtomexUpdate(ctx context.Context, data atomex.Message) error {
	switch val := data.Value.(type) {

	case atomex.OrderWebsocket:
		mm.log.Info().Int64("id", val.ID).Str("status", string(val.Status)).Str("price", val.Price.String()).Str("leave_qty", val.LeaveQty.String()).Msg("atomex order status changed")

		if err := mm.handleAtomexOrderUpdate(val); err != nil {
			return err
		}

	case atomex.Swap:
		mm.log.Info().Int64("id", val.ID).Str("user_status", string(val.User.Status)).Str("counterparty_status", string(val.CounterParty.Status)).Msg("atomex swap status changed")

		if err := mm.handleAtomexSwapUpdate(val); err != nil {
			return err
		}

		if err := mm.initiateInvolvedSwap(ctx, val); err != nil {
			return err
		}
	}

	return nil
}

func (mm *MarketMaker) sendOneByOneLimits(force bool) error {
	for i := range mm.strategies {
		if !mm.strategies[i].Is(strategy.KindOneByOne) {
			continue
		}

		quotes, err := mm.strategies[i].Quotes(nil)
		if err != nil {
			mm.log.Err(err).Msg("get strategy quotes")
			continue
		}

		for i := range quotes {
			if err := mm.sendOrder(quotes[i], force); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mm *MarketMaker) sendLimitsByTicker(tick exchange.Ticker) error {
	symbol, ok := mm.quoteProviderMeta.ToSymbols[tick.Symbol]
	if !ok {
		return errors.Errorf("unknown provider symbol: %s", tick.Symbol)
	}

	return mm.processTicker(tick, symbol)
}

func (mm *MarketMaker) processTicker(tick exchange.Ticker, symbol string) error {
	for synthSymbol, synth := range mm.synthetics {
		ticker, err := synth.Ticker(tick, mm.tickers, mm.quoteProviderMeta.ToSymbols)
		if err != nil {
			if errors.Is(err, synthetic.ErrInvalidSymbol) || errors.Is(err, synthetic.ErrUnknownTicker) {
				continue
			}
			return errors.Wrap(err, "synthetic.Ticker")
		}

		mm.tickers[ticker.Symbol] = ticker

		args := strategy.NewArgs().Ask(ticker.Ask).Bid(ticker.Bid).AskVolume(ticker.AskVolume).BidVolume(ticker.BidVolume).Symbol(synthSymbol)
		for i := range mm.strategies {
			quotes, err := mm.strategies[i].Quotes(args)
			if err != nil {
				return errors.Wrap(err, "Quotes")
			}
			for j := range quotes {
				if err := mm.sendOrder(quotes[j], false); err != nil {
					return errors.Wrap(err, "sendOrder")
				}
			}
		}
	}

	return nil
}

func (mm *MarketMaker) sendOrder(quote strategy.Quote, force bool) error {
	symbol, ok := mm.atomexMeta.ToSymbols[quote.Symbol]
	if !ok {
		return nil
	}

	symbolInfo, ok := mm.symbols[quote.Symbol]
	if !ok {
		return nil
	}

	price, _ := quote.Price.Float64()
	qty, _ := quote.Volume.Float64()
	side := mustAtomexSide(quote.Side)

	clientID := clientOrderID{
		kind:   quote.Strategy,
		symbol: quote.Symbol,
		side:   quote.Side,
		index:  time.Now().UnixNano(),
	}

	var cancelErr error
	var changed bool
	mm.orders.Range(func(cid clientOrderID, order *Order) bool {
		if clientID.Equals(cid) {
			if price != order.Price || force {
				if err := mm.atomex.CancelOrder(atomex.CancelOrderRequest{
					ID:     order.ID,
					Symbol: order.Symbol,
					Side:   order.Side,
				}); err != nil {
					cancelErr = err
				}
				mm.log.Info().Int64("id", order.ID).Msg("order cancelling...")
				changed = true
			}
			return false
		}
		return true
	})
	if cancelErr != nil {
		return errors.Wrap(cancelErr, "atomex.CancelOrder")
	}

	if !changed {
		return nil
	}

	receiver, err := mm.getReceiverWallet(symbolInfo, quote.Side)
	if err != nil {
		return errors.Wrap(err, "getReceiverWallet")
	}

	sender, err := mm.getSenderWallet(symbolInfo, quote.Side)
	if err != nil {
		return errors.Wrap(err, "getSenderWallet")
	}

	scrt, err := mm.secret(sender.Private, sender.Address, clientID.index)
	if err != nil {
		return errors.Wrap(err, "secret")
	}

	mm.secrets.Add(chain.Hex(scrt.Hash), chain.Hex(scrt.Value))

	req := atomex.NewTokenRequest(sender.Address, signers.AlgorithmEd25519Blake2b, sender.PublicKey)
	if err := req.Sign(&signers.Key{
		Public:  sender.PublicKey,
		Private: sender.Private,
	}); err != nil {
		return errors.Wrap(err, "Sign")
	}

	request := atomex.AddOrderRequest{
		ClientOrderID: clientID.String(),
		Symbol:        symbol,
		Price:         price,
		Qty:           qty,
		Side:          side,
		Type:          atomex.OrderTypeReturn,
		ProofsOfFunds: []atomex.ProofOfFunds{
			{
				Address:   sender.Address,
				Currency:  sender.Currency,
				Algorithm: req.Algorithm,
				Message:   req.Message,
				PublicKey: req.PublicKey,
				Signature: req.Signature,
				Timestamp: req.Timestamp,
			},
		},
		Requisites: &atomex.Requisites{
			BaseCurrencyContract:  symbolInfo.Base.AtomexContract,
			QuoteCurrencyContract: symbolInfo.Quote.AtomexContract,
			ReceivingAddress:      receiver.Address,
			RefundAddress:         receiver.Address,
			SecretHash:            scrt.Hash,
			LockTime:              mm.atomexMeta.Settings.LockTime,
			RewardForRedeem:       mm.atomexMeta.Settings.RewardForRedeem,
		},
	}

	if err := mm.atomex.SendOrder(request); err != nil {
		return errors.Wrap(err, "SendOrder")
	}

	order := requestToOrder(request, scrt)
	mm.orders.Store(clientID, &order)
	return nil
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

	first := sha256.Sum256(scrt)
	secretHash := sha256.Sum256(first[:])
	s.Hash = hex.EncodeToString(secretHash[:])
	return s, nil
}

func (mm *MarketMaker) findDuplicatesOrders(orders []atomex.Order) error {
	for i := range orders {
		mm.log.Info().Int64("id", orders[i].ID).Str("status", string(orders[i].Status)).Msg("find placed order")
		var clientID clientOrderID
		if err := clientID.parse(orders[i].ClientOrderID); err != nil {
			return errors.Wrap(err, "cid.parse")
		}

		found := mm.cancelOrder(clientID)
		if !found {
			internalOrder := atomexOrderToInternal(orders[i])
			mm.orders.Store(clientID, &internalOrder)
		}
	}
	return nil
}

func (mm *MarketMaker) cancelAll(ctx context.Context) (cancelErr error) {

	mm.orders.Range(func(cid clientOrderID, order *Order) bool {
		mm.log.Debug().Int64("order_id", order.ID).Str("symbol", order.Symbol).Msg("cancelling...")
		if order.ID == 0 {
			return true
		}

		cancelCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		response, err := mm.atomexAPI.CancelOrder(cancelCtx, order.ID, order.Symbol, order.Side)
		if err != nil {
			cancelErr = err
			return false
		}

		mm.swaps.Delete(chain.Hex(order.Secret.Hash))

		if response.Result {
			mm.log.Debug().Int64("order_id", order.ID).Str("symbol", order.Symbol).Msg("cancelled")

			var cid clientOrderID
			if err := cid.parse(order.ClientID); err != nil {
				cancelErr = err
				return false
			}

			mm.orders.Delete(cid)

		} else {
			mm.log.Debug().Int64("order_id", order.ID).Str("symbol", order.Symbol).Msg("not cancelled")
		}
		return true
	})

	return
}

func (mm *MarketMaker) handleAtomexSwapUpdate(swap atomex.Swap) error {
	if swap.User.Status != atomex.SwapStatusInvolved || swap.CounterParty.Status != atomex.SwapStatusInvolved {
		return nil
	}

	refundTime := swap.TimeStamp.Add(time.Duration(swap.User.Requisites.LockTime) * time.Second)
	if refundTime.Before(time.Now()) {
		mm.log.Warn().Msg("refund time has already come")
		return nil
	}

	s, err := mm.atomexSwapToInternal(swap)
	if err != nil {
		return errors.Wrap(err, "atomexSwapToInternal")
	}

	mm.swaps.Store(chain.Hex(swap.SecretHash), s)
	return nil
}

func (mm *MarketMaker) initiateInvolvedSwap(ctx context.Context, swap atomex.Swap) error {
	if swap.SecretHash == "" {
		return nil
	}
	symbol, ok := mm.atomexMeta.FromSymbols[swap.Symbol]
	if !ok {
		return errors.Errorf("unknown symbol: %s", swap.Symbol)
	}

	info, ok := mm.symbols[symbol]
	if !ok {
		return errors.Errorf("unknown symbol: %s", symbol)
	}

	var asset types.Asset
	switch swap.Side {
	case atomex.SideBuy:
		asset = info.Quote
	case atomex.SideSell:
		asset = info.Base
	}

	payOff := decimal.NewFromFloat(mm.atomexMeta.Settings.RewardForRedeem)
	refundTime := swap.TimeStamp.Add(time.Duration(swap.User.Requisites.LockTime) * time.Second)

	return mm.tracker.Initiate(ctx, chain.InitiateArgs{
		HashedSecret: chain.Hex(swap.User.Requisites.SecretHash),
		Participant:  swap.CounterParty.Requisites.ReceivingAddress,
		Contract:     asset.AtomexContract,
		TokenAddress: asset.Contract,
		Amount:       amountToInt(swap.Qty, asset.Decimals),
		PayOff:       payOff,
		RefundTime:   refundTime,
	}, asset.ChainType())
}

func (mm *MarketMaker) handleAtomexOrderUpdate(order atomex.OrderWebsocket) error {
	var cid clientOrderID
	if err := cid.parse(order.ClientOrderID); err != nil {
		return err
	}

	// handle here, because `OrderStatusFilled` or `OrderStatusPartiallyFilled` may be first and `OrderStatusPlaced` may be skipped.
	internalOrder, orderFound := mm.orders.Load(cid)
	if orderFound {
		internalOrder.Status = string(order.Status)
		internalOrder.ID = order.ID
	}

	switch order.Status {
	case atomex.OrderStatusCanceled, atomex.OrderStatusRejected:
		mm.orders.Delete(cid)
		if internalOrder != nil {
			mm.swaps.Delete(chain.Hex(internalOrder.Secret.Hash))
		}

	case atomex.OrderStatusPartiallyFilled, atomex.OrderStatusFilled:
		if err := mm.sendOneByOneLimits(true); err != nil {
			return err
		}
	case atomex.OrderStatusPending: // do not handle. it's internal atomex status.
	case atomex.OrderStatusPlaced:
	}
	return nil
}

func (mm *MarketMaker) cancelOrder(clientID clientOrderID) bool {
	var found bool
	mm.orders.Range(func(cid clientOrderID, order *Order) bool {
		if clientID.Equals(cid) {
			mm.log.Warn().Int64("id", order.ID).Msg("order will be cancelled.")
			if err := mm.atomex.CancelOrder(atomex.CancelOrderRequest{
				ID:     order.ID,
				Symbol: order.Symbol,
				Side:   order.Side,
			}); err != nil {
				mm.log.Error().Err(err).Msg("order cancelling")
			} else {
				mm.log.Info().Int64("id", order.ID).Msg("order cancelling...")
				mm.secrets.Delete(chain.Hex(order.Secret.Hash))
			}
			mm.orders.Delete(cid)
			found = true
			return false
		}
		return true
	})
	return found
}
