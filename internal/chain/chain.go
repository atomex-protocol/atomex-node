package chain

import (
	"context"
	"io"
	"math/big"
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
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Initiate(ctx context.Context, args InitiateArgs) error
	Redeem(ctx context.Context, hashedSecret, secret Hex, contract string) error
	Refund(ctx context.Context, hashedSecret Hex, contract string) error
	Restore(ctx context.Context) error
	Wallet() Wallet

	Events() <-chan Event
	Operations() <-chan Operation
}

// InitiateArgs -
type InitiateArgs struct {
	HashedSecret Hex
	Participant  string
	Contract     string
	TokenAddress string
	Amount       decimal.Decimal
	PayOff       decimal.Decimal
	RefundTime   time.Time
}

// Wallet -
type Wallet struct {
	Address   string
	PublicKey []byte
	Private   []byte
}

// Event -
type Event interface {
	Level() uint64
	HashedSecret() Hex
	ChainType() ChainType
	Contract() string
}

// RedeemEvent -
type RedeemEvent struct {
	HashedSecretHex Hex
	ContractAddress string
	Chain           ChainType
	BlockNumber     uint64
	Secret          Hex
}

// Level -
func (e RedeemEvent) Level() uint64 {
	return e.BlockNumber
}

// Contract -
func (e RedeemEvent) Contract() string {
	return e.ContractAddress
}

// ChainType -
func (e RedeemEvent) ChainType() ChainType {
	return e.Chain
}

// HashedSecret -
func (e RedeemEvent) HashedSecret() Hex {
	return e.HashedSecretHex
}

// RefundEvent -
type RefundEvent struct {
	HashedSecretHex Hex
	ContractAddress string
	Chain           ChainType
	BlockNumber     uint64
}

// Level -
func (e RefundEvent) Level() uint64 {
	return e.BlockNumber
}

// Contract -
func (e RefundEvent) Contract() string {
	return e.ContractAddress
}

// ChainType -
func (e RefundEvent) ChainType() ChainType {
	return e.Chain
}

// HashedSecret -
func (e RefundEvent) HashedSecret() Hex {
	return e.HashedSecretHex
}

// InitEvent -
type InitEvent struct {
	HashedSecretHex Hex
	ContractAddress string
	Chain           ChainType
	BlockNumber     uint64
	Initiator       string
	Participant     string
	Amount          decimal.Decimal
	PayOff          decimal.Decimal
	RefundTime      time.Time
}

// Level -
func (e *InitEvent) Level() uint64 {
	return e.BlockNumber
}

// Contract -
func (e *InitEvent) Contract() string {
	return e.ContractAddress
}

// ChainType -
func (e *InitEvent) ChainType() ChainType {
	return e.Chain
}

// HashedSecret -
func (e *InitEvent) HashedSecret() Hex {
	return e.HashedSecretHex
}

// SetPayOff -
func (e *InitEvent) SetPayOff(payoff *big.Int, minPayoff decimal.Decimal) error {
	if payoff == nil { // If payoff is empty or not supported.
		e.PayOff = decimal.Zero
		return nil
	}

	payoffDecimal := decimal.NewFromBigInt(payoff, 0)
	if minPayoff.Cmp(payoffDecimal) > 0 {
		return ErrMinPayoff
	}
	e.PayOff = payoffDecimal
	return nil
}

// RestoredEvent -
type RestoredEvent struct {
	Chain ChainType
}

// Level -
func (e RestoredEvent) Level() uint64 {
	return 0
}

// Contract -
func (e RestoredEvent) Contract() string {
	return ""
}

// ChainType -
func (e RestoredEvent) ChainType() ChainType {
	return e.Chain
}

// HashedSecret -
func (e RestoredEvent) HashedSecret() Hex {
	return ""
}

// ChainType -
type ChainType int

// chain types
const (
	ChainTypeUnknown ChainType = iota
	ChainTypeTezos
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

// ByLevel -
type ByLevel []Event

// Len -
func (a ByLevel) Len() int { return len(a) }

// Less -
func (a ByLevel) Less(i, j int) bool { return a[i].Level() < a[j].Level() }

// Swap -
func (a ByLevel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
