package tools

import (
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// Swap -
type Swap struct {
	HashedSecret    chain.Hex
	Secret          chain.Hex
	Status          Status
	RefundTime      time.Time
	RewardForRedeem decimal.Decimal
	Initiator       Leg
	Acceptor        Leg
	Symbol          types.Symbol
}

// NewSwap -
func NewSwap(event chain.Event) *Swap {
	return &Swap{
		HashedSecret: event.HashedSecret(),
		Status:       StatusEmpty,
	}
}

// Log -
func (swap *Swap) Log(logger *zerolog.Event) *zerolog.Event {
	return logger.Str("hashed_secret", swap.HashedSecret.String()).
		Str("status", swap.Status.String()).
		Str("initiator_chain", swap.Initiator.ChainType.String()).
		Str("acceptor_chain", swap.Acceptor.ChainType.String())
}

// IsUnknown -
func (swap *Swap) IsUnknown() bool {
	return swap.Acceptor.ChainType == chain.ChainTypeUnknown || swap.Initiator.ChainType == chain.ChainTypeUnknown
}

func (swap *Swap) fromInitEvent(event chain.InitEvent) {
	if swap.HashedSecret != event.HashedSecret() {
		return
	}

	swap.RefundTime = event.RefundTime
	swap.RewardForRedeem = event.PayOff

	switch swap.Status {
	case StatusEmpty:
		swap.Initiator = Leg{
			ChainType: event.Chain,
			Address:   event.Initiator,
			Contract:  event.ContractAddress,
			Status:    StatusInitiated,
		}
		swap.Acceptor = Leg{
			Address: event.Participant,
			Status:  StatusEmpty,
		}
		swap.Status = StatusInitiatedOnce
	case StatusInitiatedOnce:
		swap.Acceptor.ChainType = event.Chain
		swap.Acceptor.Contract = event.ContractAddress
		swap.Acceptor.Status = StatusInitiated
		swap.Status = StatusInitiated
	}
}

func (swap *Swap) fromRedeemEvent(event chain.RedeemEvent) {
	if swap.HashedSecret != event.HashedSecret() {
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

	if swap.Acceptor.Contract == event.ContractAddress && swap.Acceptor.ChainType == event.Chain {
		swap.Acceptor.Status = StatusRedeemed
	}
	if swap.Initiator.Contract == event.ContractAddress && swap.Initiator.ChainType == event.Chain {
		swap.Initiator.Status = StatusRedeemed
	}
}

func (swap *Swap) fromRefundEvent(event chain.RefundEvent) {
	if swap.HashedSecret != event.HashedSecret() {
		return
	}
	switch swap.Status {
	case StatusEmpty, StatusInitiatedOnce, StatusInitiated:
		swap.Status = StatusRefundedOnce
	case StatusRefundedOnce:
		swap.Status = StatusRefunded
	}

	if swap.Acceptor.Contract == event.ContractAddress && swap.Acceptor.ChainType == event.Chain {
		swap.Acceptor.Status = StatusRefunded
	}
	if swap.Initiator.Contract == event.ContractAddress && swap.Initiator.ChainType == event.Chain {
		swap.Initiator.Status = StatusRefunded
	}
}

// Leg -
type Leg struct {
	ChainType chain.ChainType
	Address   string
	Contract  string
	Status    Status
}

// IsFinished -
func (leg Leg) IsFinished() bool {
	return leg.Status == StatusRedeemed || leg.Status == StatusRefunded
}

// Merge -
func (leg *Leg) Merge(another Leg) {
	if leg.ChainType != another.ChainType {
		return
	}
	if leg.Status < another.Status {
		leg.Status = another.Status
	}
	if leg.Contract == "" {
		leg.Contract = another.Contract
	}
	if leg.Address == "" {
		leg.Address = another.Address
	}
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
