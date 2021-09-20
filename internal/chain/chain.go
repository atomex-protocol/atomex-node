package chain

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// errors
var (
	ErrMinPayoff = errors.New("payoff threashold")
)

// Chain -
type Chain interface {
	io.Closer
	Run() error
	Redeem(hashedSecret, secret Hex, contract string) error
	Refund(hashedSecret Hex, contract string) error
	Restore() error

	InitEvents() <-chan InitEvent
	RedeemEvents() <-chan RedeemEvent
	RefundEvents() <-chan RefundEvent
	Operations() <-chan Operation
}

// RedeemEvent -
type RedeemEvent struct {
	Event
	Secret Hex
}

// RefundEvent -
type RefundEvent struct {
	Event
}

// InitEvent -
type InitEvent struct {
	Event
	Initiator   string
	Participant string
	Amount      decimal.Decimal
	PayOff      decimal.Decimal
	RefundTime  time.Time
}

// Event -
type Event struct {
	HashedSecret Hex
	Contract     string
	Chain        ChainType
}

// SetAmountFromString -
func (event *InitEvent) SetAmountFromString(amount string) error {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return err
	}
	event.Amount = amountDecimal
	return nil
}

// SetAmountFromString -
func (event *InitEvent) SetPayOff(payoff string, minPayoff decimal.Decimal) error {
	payoffDecimal, err := decimal.NewFromString(payoff)
	if err != nil {
		return err
	}
	if minPayoff.Cmp(payoffDecimal) > 0 {
		return ErrMinPayoff
	}
	event.PayOff = payoffDecimal
	return nil
}

// ChainType -
type ChainType int

// chain types
const (
	ChainTypeTezos ChainType = iota + 1
	ChainTypeEthereum
)

// String -
func (c ChainType) String() string {
	switch c {
	case ChainTypeTezos:
		return "tezos"
	case ChainTypeEthereum:
		return "ethereum"
	default:
		return "unknown"
	}
}
