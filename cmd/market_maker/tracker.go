package main

import (
	"context"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/pkg/errors"
)

func (mm *MarketMaker) listenTracker(ctx context.Context) {
	defer mm.wg.Done()

	for {
		select {
		case <-ctx.Done():
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
			if len(current.Secret) == 0 {
				current.Secret = swap.Secret
			}

			switch current.Status {
			case tools.StatusInitiated:
				if current.IsUnknown() {
					continue
				}

				mm.log.Info().Str("hashed_secret", current.HashedSecret.String()).Msg("swap is initiated. redeeming...")

				if err := mm.restoreSecretFromTrackerAtomex(ctx, current); err != nil {
					mm.log.Err(err).Msg("restoreSecretFromTrackerAtomex")
					continue
				}

				if len(current.Secret) == 0 {
					mm.log.Error().Str("hashed_secret", current.HashedSecret.String()).Msg("empty secret before redeem")
					continue
				}

				if err := mm.tracker.Redeem(ctx, *current, current.Acceptor); err != nil {
					mm.log.Err(err).Msg("tracker.Redeem")
					continue
				}

			case tools.StatusRefundedOnce:
				if current.IsUnknown() {
					continue
				}

				mm.log.Info().Str("hashed_secret", current.HashedSecret.String()).Msg("counterparty refunded swap. refunding...")

				if err := mm.restoreSecretFromTrackerAtomex(ctx, current); err != nil {
					mm.log.Err(err).Msg("restoreSecretFromTrackerAtomex")
					continue
				}
				if len(current.Secret) == 0 {
					mm.log.Error().Str("hashed_secret", current.HashedSecret.String()).Msg("empty secret before refund")
					continue
				}

				if err := mm.tracker.Refund(ctx, *current, current.Initiator); err != nil {
					mm.log.Err(err).Msg("tracker.Refund")
					continue
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

func (mm *MarketMaker) restoreSecretFromTrackerAtomex(ctx context.Context, swap *tools.Swap) error {
	if !swap.Secret.IsEmpty() {
		return nil
	}

	for i := range mm.activeSwaps {
		if mm.activeSwaps[i].SecretHash != swap.HashedSecret.String() {
			continue
		}

		if err := mm.restoreSecretForAtomexSwap(ctx, mm.activeSwaps[i], swap); err != nil {
			return errors.Wrap(err, "restoreSecretForAtomexSwap")
		}

		break
	}

	return nil
}
