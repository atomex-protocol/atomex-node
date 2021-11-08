package main

import (
	"sync"

	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
)

// Order -
type Order struct {
	ClientID string
	ID       int64
	Symbol   string
	Status   string
	Price    float64
	Qty      float64
	Side     atomex.Side
	Type     atomex.OrderType
	Secret   secret
}

// OrdersMap
type OrdersMap struct {
	mx sync.RWMutex
	m  map[clientOrderID]*Order
}

// NewOrdersMap -
func NewOrdersMap() *OrdersMap {
	return &OrdersMap{
		m: make(map[clientOrderID]*Order),
	}
}

// Load -
func (orders *OrdersMap) Load(key clientOrderID) (*Order, bool) {
	orders.mx.RLock()
	val, ok := orders.m[key]
	orders.mx.RUnlock()
	return val, ok
}

// Set -
func (orders *OrdersMap) Store(cid clientOrderID, order *Order) {
	orders.mx.Lock()
	orders.m[cid] = order
	orders.mx.Unlock()
}

// Delete -
func (orders *OrdersMap) Delete(cid clientOrderID) {
	orders.mx.Lock()
	delete(orders.m, cid)
	orders.mx.Unlock()
}

// Range -
func (orders *OrdersMap) Range(handler func(cid clientOrderID, order *Order) bool) {
	orders.mx.RLock()

	for key, value := range orders.m {
		orders.mx.RUnlock()
		if !handler(key, value) {
			return
		}
		orders.mx.RLock()
	}

	orders.mx.RUnlock()
}

// SwapsMap
type SwapsMap struct {
	mx sync.RWMutex
	m  map[chain.Hex]*tools.Swap
}

// NewSwapsMap -
func NewSwapsMap() *SwapsMap {
	return &SwapsMap{
		m: make(map[chain.Hex]*tools.Swap),
	}
}

// Load -
func (swaps *SwapsMap) Load(hashedSecret chain.Hex) (*tools.Swap, bool) {
	swaps.mx.RLock()
	val, ok := swaps.m[hashedSecret]
	swaps.mx.RUnlock()
	return val, ok
}

// Set -
func (swaps *SwapsMap) Store(hashedSecret chain.Hex, swap *tools.Swap) {
	swaps.mx.Lock()
	swaps.m[hashedSecret] = swap
	swaps.mx.Unlock()
}

// Load -
func (swaps *SwapsMap) LoadOrStore(hashedSecret chain.Hex, swap *tools.Swap) *tools.Swap {
	val, ok := swaps.Load(hashedSecret)
	if !ok {
		return val
	}

	swaps.Store(hashedSecret, swap)
	return swap
}

// Delete -
func (swaps *SwapsMap) Delete(hashedSecret chain.Hex) {
	swaps.mx.Lock()
	delete(swaps.m, hashedSecret)
	swaps.mx.Unlock()
}

// Range -
func (swaps *SwapsMap) Range(handler func(hashedSecret chain.Hex, swap *tools.Swap) bool) {
	swaps.mx.RLock()

	for key, value := range swaps.m {
		swaps.mx.RUnlock()
		if !handler(key, value) {
			return
		}
		swaps.mx.RLock()
	}

	swaps.mx.RUnlock()
}
