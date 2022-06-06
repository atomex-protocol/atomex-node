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
