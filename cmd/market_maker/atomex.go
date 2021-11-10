package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/synthetic"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (mm *MarketMaker) listenAtomex() {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.stop:
			return

		case data := <-mm.atomex.Listen():
			if err := mm.handleAtomexUpdate(data); err != nil {
				mm.log.Err(err).Msg("handleAtomexUpdate")
				continue
			}

		case err := <-mm.atomex.Errors():
			mm.log.Err(err).Msg("atomex error")
		}
	}
}

func (mm *MarketMaker) handleAtomexUpdate(data atomex.Message) error {
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
	}

	return nil
}

func (mm *MarketMaker) sendOneByOneLimits() error {
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
			if err := mm.sendOrder(quotes[i]); err != nil {
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
				if err := mm.sendOrder(quotes[j]); err != nil {
					return errors.Wrap(err, "sendOrder")
				}
			}
		}
	}

	return nil
}

func (mm *MarketMaker) sendOrder(quote strategy.Quote) error {
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
		index:  atomic.AddUint64(&mm.ordersCounter, 1),
	}

	var cancelErr error
	mm.orders.Range(func(cid clientOrderID, order *Order) bool {
		if cid.kind == clientID.kind && cid.side == clientID.side && cid.symbol == clientID.symbol {
			if err := mm.atomex.CancelOrder(atomex.CancelOrderRequest{
				ID:     order.ID,
				Symbol: order.Symbol,
				Side:   order.Side,
			}); err != nil {
				cancelErr = err
				return false
			}
		}
		return true
	})
	if cancelErr != nil {
		return errors.Wrap(cancelErr, "atomex.CancelOrder")
	}

	receiver, err := mm.getReceiverWallet(symbolInfo, quote.Side)
	if err != nil {
		return errors.Wrap(err, "getReceiverWallet")
	}

	sender, err := mm.getSenderWallet(symbolInfo, quote.Side)
	if err != nil {
		return errors.Wrap(err, "getSenderWallet")
	}

	scrt, err := mm.secret(sender.Private, sender.Address, time.Now().UnixNano())
	if err != nil {
		return errors.Wrap(err, "secret")
	}

	request := atomex.AddOrderRequest{
		ClientOrderID: clientID.String(),
		Symbol:        symbol,
		Price:         price,
		Qty:           qty,
		Side:          side,
		Type:          atomex.OrderTypeReturn,
		Requisites: &atomex.Requisites{
			BaseCurrencyContract:  symbolInfo.Base.AtomexContract,
			QuoteCurrencyContract: symbolInfo.Quote.AtomexContract,
			SecretHash:            scrt.Hash,
			ReceivingAddress:      receiver.Address,
			RefundAddress:         receiver.Address,
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

	secretHash := sha256.Sum256(scrt)
	s.Hash = hex.EncodeToString(secretHash[:])
	return s, nil
}

func (mm *MarketMaker) findDuplicatesOrders(orders []atomex.Order) error {
	for i := range orders {
		mm.log.Info().Int64("id", orders[i].ID).Str("status", string(orders[i].Status)).Msg("atomex order status changed")
		var cid clientOrderID
		if err := cid.parse(orders[i].ClientOrderID); err != nil {
			return errors.Wrap(err, "cid.parse")
		}

		if order, ok := mm.orders.Load(cid); ok {
			if order.ID == orders[i].ID {
				continue
			}

			mm.log.Warn().Int64("id", order.ID).Int64("second_id", orders[i].ID).Msg("found order duplicate. it will be cancelled.")
			if err := mm.atomex.CancelOrder(atomex.CancelOrderRequest{
				ID:     orders[i].ID,
				Side:   orders[i].Side,
				Symbol: orders[i].Symbol,
			}); err != nil {
				return errors.Wrap(err, "atomex.CancelOrder")
			}
		} else {
			internalOrder := atomexOrderToInternal(orders[i])
			mm.orders.Store(cid, &internalOrder)
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

	s, err := mm.atomexSwapToInternal(swap)
	if err != nil {
		return errors.Wrap(err, "atomexSwapToInternal")
	}

	mm.swaps.Store(chain.Hex(swap.SecretHash), s)

	var asset types.Asset
	switch swap.Side {
	case atomex.SideBuy:
		asset = s.Symbol.Quote
	case atomex.SideSell:
		asset = s.Symbol.Base
	}

	payOff := decimal.NewFromFloat(mm.atomexMeta.Settings.RewardForRedeem)
	refundTime := swap.TimeStamp.Add(time.Duration(swap.User.Requisites.LockTime) * time.Second)

	if err := mm.tracker.Initiate(chain.InitiateArgs{
		HashedSecret: chain.Hex(swap.User.Requisites.SecretHash),
		Participant:  swap.CounterParty.Requisites.ReceivingAddress,
		Contract:     asset.Contract,
		TokenAddress: asset.Contract,
		Amount:       swap.Qty,
		PayOff:       payOff,
		RefundTime:   refundTime,
	}, asset.ChainType()); err != nil {
		return errors.Wrap(err, "Initiate")
	}

	return nil
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

	case atomex.OrderStatusFilled, atomex.OrderStatusPartiallyFilled: // do not handle. it's because it's handled in `handleAtomexSwapUpdate`
	case atomex.OrderStatusPending: // do not handle. it's internal atomex status.
	case atomex.OrderStatusPlaced: // do not handle, because it's handled above
	}
	return nil
}
