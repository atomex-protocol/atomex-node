package main

import (
	"math/rand"
	"time"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateClientOrderID() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
