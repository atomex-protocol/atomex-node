package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/internal/atomex"
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/pkg/errors"
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

func (orders *OrdersMap) Count() int {
	orders.mx.RLock()
	count := len(orders.m)
	orders.mx.RUnlock()
	return count
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

// Store -
func (swaps *SwapsMap) Store(hashedSecret chain.Hex, swap *tools.Swap) {
	swaps.mx.Lock()
	swaps.m[hashedSecret] = swap
	swaps.mx.Unlock()
}

// LoadOrStore -
func (swaps *SwapsMap) LoadOrStore(hashedSecret chain.Hex, swap *tools.Swap) *tools.Swap {
	val, ok := swaps.Load(hashedSecret)
	if ok {
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

type clientOrderID struct {
	kind   strategy.Kind
	symbol string
	side   strategy.Side
	index  int64
}

func (c clientOrderID) String() string {
	return fmt.Sprintf("%d%d%d%s", strategyKindToInt(c.kind), c.side, c.index, c.symbol)
}

func (c *clientOrderID) parse(str string) error {
	if len(str) < 22 {
		return errors.Errorf("invalid client order id '%s'", str)
	}

	kind, err := strconv.ParseInt(str[0:1], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid client order id '%s'", str)
	}
	c.kind = intToStrategyKind(kind)

	side, err := strconv.ParseInt(str[1:2], 10, 32)
	if err != nil {
		return errors.Wrapf(err, "invalid client order id '%s'", str)
	}
	c.side = strategy.Side(side)

	index, err := strconv.ParseInt(str[2:21], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalid client order id '%s'", str)
	}
	c.index = index

	c.symbol = str[21:]

	return nil
}

// Equals -
func (c clientOrderID) Equals(clientID clientOrderID) bool {
	return c.kind == clientID.kind && c.side == clientID.side && c.symbol == clientID.symbol
}

func strategyKindToInt(kind strategy.Kind) int {
	switch kind {
	case strategy.KindFollow:
		return 1
	case strategy.KindOneByOne:
		return 2
	case strategy.KindVolatility:
		return 3
	}
	return 0
}

func intToStrategyKind(kind int64) strategy.Kind {
	switch kind {
	case 1:
		return strategy.KindFollow
	case 2:
		return strategy.KindOneByOne
	case 3:
		return strategy.KindVolatility
	}
	return strategy.KindUnknown
}

// Wallet -
type Wallet struct {
	chain.Wallet
	Currency string
}

func newWallet(wallet chain.Wallet, asset types.Asset) Wallet {
	return Wallet{wallet, asset.Name}
}

// Secrets -
type Secrets struct {
	m  map[chain.Hex]chain.Hex
	mx sync.RWMutex
}

// NewSecrets -
func NewSecrets() *Secrets {
	return &Secrets{
		m: make(map[chain.Hex]chain.Hex),
	}
}

// Get -
func (s *Secrets) Get(hash chain.Hex) (chain.Hex, bool) {
	s.mx.RLock()
	secret, ok := s.m[hash]
	s.mx.RUnlock()
	return secret, ok
}

// Add -
func (s *Secrets) Add(hash, secret chain.Hex) {
	s.mx.Lock()
	s.m[hash] = secret
	s.mx.Unlock()
}

// Delete -
func (s *Secrets) Delete(hash chain.Hex) {
	s.mx.Lock()
	delete(s.m, hash)
	s.mx.Unlock()
}
