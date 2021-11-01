package main

import "github.com/atomex-protocol/watch_tower/internal/atomex"

// Order -
type Order struct {
	ClientID string
	ID       string
	Symbol   string
	Price    float64
	Qty      float64
	Side     atomex.Side
	Type     atomex.OrderType
	Secret   secret
}
