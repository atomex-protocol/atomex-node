package main

import (
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/shopspring/decimal"
)

func mustAtomexSide(side strategy.Side) atomex.Side {
	switch side {
	case strategy.Ask:
		return atomex.SideSell
	case strategy.Bid:
		return atomex.SideBuy
	}
	return ""
}

func requestToOrder(request atomex.AddOrderRequest, secret secret) Order {
	return Order{
		ClientID: request.ClientOrderID,
		Symbol:   request.Symbol,
		Side:     request.Side,
		Secret:   secret,
		Price:    request.Price,
		Qty:      request.Qty,
		Type:     request.Type,
	}
}

func atomexOrderToInternal(order atomex.Order) Order {
	price, _ := order.Price.Float64()
	qty, _ := order.Qty.Float64()

	var scrt, hashedScrt string
	if len(order.Swaps) > 0 {
		scrt = order.Swaps[0].Secret
		hashedScrt = order.Swaps[0].SecretHash
	}

	return Order{
		ID:       order.ID,
		ClientID: order.ClientOrderID,
		Symbol:   order.Symbol,
		Side:     order.Side,
		Secret: secret{
			Value: scrt,
			Hash:  hashedScrt,
		},
		Price: price,
		Qty:   qty,
		Type:  order.Type,
	}
}

func atomexStatusToInternal(status atomex.SwapStatus) tools.Status {
	switch status {
	case atomex.SwapStatusInitiated:
		return tools.StatusInitiated
	case atomex.SwapStatusRedeemed:
		return tools.StatusRedeemed
	case atomex.SwapStatusRefunded:
		return tools.StatusRefunded
	}
	return tools.StatusEmpty
}

func amountToInt(value decimal.Decimal, decimals int) decimal.Decimal {
	mux := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))
	return value.Mul(mux).Round(0)
}
