package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/pkg/errors"
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

type clientOrderID struct {
	kind   strategy.Kind
	symbol string
	side   strategy.Side
	index  int
}

func (c clientOrderID) String() string {
	return fmt.Sprintf("%s|%s|%d|%d", c.kind, c.symbol, c.side, c.index)
}

func (c *clientOrderID) parse(str string) error {
	parts := strings.Split(str, "|")
	if len(parts) != 4 {
		return errors.Errorf("invalid client order id '%s'", str)
	}

	c.kind = strategy.Kind(parts[0])
	c.symbol = parts[1]

	side, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return errors.Wrapf(err, "invalid client order id '%s'", str)
	}
	c.side = strategy.Side(side)

	index, err := strconv.ParseInt(parts[3], 10, 32)
	if err != nil {
		return errors.Wrapf(err, "invalid client order id '%s'", str)
	}
	c.index = int(index)

	return nil
}
