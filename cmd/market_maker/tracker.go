package main

import (
	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
)

func (mm *MarketMaker) listenTracker() {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.stop:
			return

		case operation := <-mm.tracker.Operations():
			mm.log.Info().Str("hash", operation.Hash).Str("chain", operation.ChainType.String()).Str("status", operation.Status.String()).Msg("operation's status changed")

			switch operation.Status {
			case chain.Applied, chain.Pending: // do not handle. it's handled below in `StatusChanged`
			case chain.Failed:
				// TODO: think about re-send reasons
			default:
				mm.log.Warn().Msgf("unknown operation status: %s", operation.Status.String())
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
