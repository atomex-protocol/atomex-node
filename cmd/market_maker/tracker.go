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
			current.Acceptor.Merge(swap.Acceptor)
			current.Initiator.Merge(swap.Initiator)

			if current.Secret.IsEmpty() {
				current.Secret = swap.Secret
			}

			switch swap.Status {
			case tools.StatusInitiated:
				if swap.IsUnknown() {
					continue
				}

				mm.log.Info().Str("hashed_secret", swap.HashedSecret.String()).Msg("swap is initiated. redeeming...")

				if swap.Secret.IsEmpty() {
					secret, ok := mm.secrets.Get(swap.HashedSecret)
					if !ok {
						if err := mm.restoreSecretFromTrackerAtomex(ctx, &swap); err != nil {
							mm.log.Err(err).Msg("restoreSecretFromTrackerAtomex")
							continue
						}
						if swap.Secret.IsEmpty() {
							mm.log.Error().Str("hashed_secret", swap.HashedSecret.String()).Msg("empty secret before redeem")
							continue
						}
					} else {
						swap.Secret = secret
					}
				}

				if err := mm.tracker.Redeem(ctx, swap, swap.Acceptor); err != nil {
					mm.log.Err(err).Msg("tracker.Redeem")
				}
				continue

			case tools.StatusRefundedOnce:
				if swap.IsUnknown() {
					continue
				}

				mm.log.Info().Str("hashed_secret", swap.HashedSecret.String()).Msg("counterparty refunded swap. refunding...")

				if swap.Secret.IsEmpty() {
					secret, ok := mm.secrets.Get(swap.HashedSecret)
					if !ok {
						if err := mm.restoreSecretFromTrackerAtomex(ctx, &swap); err != nil {
							mm.log.Err(err).Msg("restoreSecretFromTrackerAtomex")
							continue
						}
						if swap.Secret.IsEmpty() {
							mm.log.Error().Str("hashed_secret", swap.HashedSecret.String()).Msg("empty secret before refund")
							continue
						}
					} else {
						swap.Secret = secret
					}
				}

				if err := mm.tracker.Refund(ctx, swap, swap.Initiator); err != nil {
					mm.log.Err(err).Msg("tracker.Refund")
				}

			case tools.StatusRefunded, tools.StatusRedeemed:
				mm.swaps.Delete(current.HashedSecret)
				mm.secrets.Delete(current.HashedSecret)
			}

		case <-mm.tracker.Restored():
			if err := mm.initialize(ctx); err != nil {
				mm.log.Err(err).Msg("initialize")
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
