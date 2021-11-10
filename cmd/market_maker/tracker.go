package main

import (
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/pkg/errors"
)

func (mm *MarketMaker) listenTracker() {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.stop:
			return

		case operation := <-mm.tracker.Operations():
			if err := mm.handleOperationFromChain(operation); err != nil {
				mm.log.Err(err).Msg("handleOperationFromChain")
			}

		case swap := <-mm.tracker.StatusChanged():
			swap.Log(mm.log.Info()).Msg("swap's status changed")

			current := mm.swaps.LoadOrStore(swap.HashedSecret, &swap)
			current.Status = swap.Status
			current.Acceptor.Status = swap.Acceptor.Status
			current.Initiator.Status = swap.Initiator.Status
			current.Secret = swap.Secret

			switch current.Status {
			case tools.StatusInitiated:
				mm.log.Info().Str("hashed_secret", current.HashedSecret.String()).Msg("swap is waiting redeem operation from watch tower")
			case tools.StatusRedeemedOnce:
				if current.Initiator.Status != tools.StatusRedeemed {
					if err := mm.tracker.Redeem(*current, current.Initiator); err != nil {
						mm.log.Err(err).Msg("tracker.Redeem")
						continue
					}
				}
			case tools.StatusRefundedOnce:
				if current.Initiator.Status != tools.StatusRefunded {
					if err := mm.tracker.Refund(*current, current.Initiator); err != nil {
						mm.log.Err(err).Msg("tracker.Refund")
						continue
					}
				}
			case tools.StatusRefunded, tools.StatusRedeemed:
				mm.swaps.Delete(current.HashedSecret)
			}

		}
	}
}

func (mm *MarketMaker) handleOperationFromChain(operation chain.Operation) error {
	mm.log.Info().Str("hash", operation.Hash).Str("chain", operation.ChainType.String()).Str("status", operation.Status.String()).Msg("operation's status changed")
	id := tools.OperationID{
		Hash:  operation.Hash,
		Chain: operation.ChainType,
	}

	switch operation.Status {
	case chain.Pending:
		mm.operations[id] = operation
		mm.log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", operation.HashedSecret.String()).Msg("transaction")

	case chain.Applied:
		if old, ok := mm.operations[id]; ok {
			mm.log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", old.HashedSecret.String()).Msg("transaction")
			delete(mm.operations, id)
		}
	case chain.Failed:
		if old, ok := mm.operations[id]; ok {
			mm.log.Info().Str("blockchain", operation.ChainType.String()).Str("hash", operation.Hash).Str("status", operation.Status.String()).Str("hashed_secret", old.HashedSecret.String()).Msg("transaction")
			delete(mm.operations, id)

			// TODO: resend here if needed
		}
	default:
		return errors.Errorf("unknown operation status: %s", operation.Status.String())
	}
	return nil
}
