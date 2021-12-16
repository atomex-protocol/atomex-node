package main

import "context"

func (mm *MarketMaker) listenProvider(ctx context.Context) {
	defer mm.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case tick := <-mm.provider.Tickers():
			mm.log.Debug().Str("ask", tick.Ask.String()).Str("bid", tick.Bid.String()).Str("symbol", tick.Symbol).Msg("quote provider's tick")

			if err := mm.sendLimitsByTicker(tick); err != nil {
				mm.log.Err(err).Msg("sendLimitsByTicker")
				continue
			}

		}
	}
}
