package main

import (
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
)

// Swap -
type Swap struct {
	HashedSecret chain.Hex
	Secret       chain.Hex
	Contract     string
	Status       Status
	RefundTime   time.Time
	Initiator    Leg
	Acceptor     Leg
	RetryCount   uint
}

// NewSwap -
func NewSwap(event chain.Event) *Swap {
	return &Swap{
		HashedSecret: event.HashedSecret,
		Contract:     event.Contract,
		Status:       StatusEmpty,
	}
}

// FromInitEvent -
func (swap *Swap) FromInitEvent(event chain.InitEvent) {
	if swap.HashedSecret != event.HashedSecret {
		return
	}

	swap.RefundTime = event.RefundTime

	switch swap.Status {
	case StatusEmpty:
		swap.Initiator = Leg{
			ChainType: event.Chain,
			Address:   event.Initiator,
		}
		swap.Acceptor = Leg{
			Address: event.Initiator,
		}
		swap.Status = StatusInitiatedOnce
	case StatusInitiatedOnce:
		swap.Acceptor.ChainType = event.Chain
		swap.Status = StatusInitiated
	}
}

// FromRedeemEvent -
func (swap *Swap) FromRedeemEvent(event chain.RedeemEvent) {
	if swap.HashedSecret != event.HashedSecret {
		return
	}
	if swap.Secret == "" {
		swap.Secret = event.Secret
	}

	switch swap.Status {
	case StatusEmpty, StatusInitiatedOnce, StatusInitiated:
		swap.Status = StatusRedeemedOnce
	case StatusRedeemedOnce:
		swap.Status = StatusRedeemed
	}
}

// FromRefundEvent -
func (swap *Swap) FromRefundEvent(event chain.RefundEvent) {
	if swap.HashedSecret != event.HashedSecret {
		return
	}
	switch swap.Status {
	case StatusEmpty, StatusInitiatedOnce, StatusInitiated:
		swap.Status = StatusRefundedOnce
	case StatusRefundedOnce:
		swap.Status = StatusRefunded
	}
}

// Leg -
type Leg struct {
	ChainType chain.ChainType
	Address   string
}

// Status -
type Status int

// String -
func (s Status) String() string {
	switch s {
	case StatusEmpty:
		return "new"
	case StatusInitiatedOnce:
		return "initiated_once"
	case StatusInitiated:
		return "initiated"
	case StatusRedeemedOnce:
		return "redeemed_once"
	case StatusRedeemed:
		return "redeemed"
	case StatusRefundedOnce:
		return "refunded_once"
	case StatusRefunded:
		return "refunded"
	default:
		return "unknown"
	}
}

// statuses
const (
	StatusEmpty Status = iota
	StatusInitiatedOnce
	StatusInitiated
	StatusRedeemedOnce
	StatusRedeemed
	StatusRefundedOnce
	StatusRefunded
)

// OperationID -
type OperationID struct {
	Hash  string
	Chain chain.ChainType
}
