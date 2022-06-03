package main

import (
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
)

// Swap -
type Swap struct {
	tools.Swap
	RetryCount uint
}

// Leg -
func (swap *Swap) Leg() *tools.Leg {
	if swap.Acceptor.ChainType == chain.ChainTypeUnknown || swap.Initiator.ChainType == chain.ChainTypeUnknown {
		return nil
	}

	if swap.Acceptor.IsFinished() && swap.Initiator.IsFinished() {
		return nil
	}

	switch {
	case swap.Acceptor.Status > swap.Initiator.Status:
		return &swap.Initiator
	case swap.Acceptor.Status < swap.Initiator.Status:
		return &swap.Acceptor
	}
	return nil
}

func (swap *Swap) merge(update tools.Swap) {
	if update.HashedSecret.String() != swap.HashedSecret.String() || update.Status < swap.Status {
		return
	}

	swap.Status = update.Status
	swap.Acceptor.Merge(update.Acceptor)
	swap.Initiator.Merge(update.Initiator)

	if swap.Secret.IsEmpty() {
		swap.Secret = update.Secret
	}
}
