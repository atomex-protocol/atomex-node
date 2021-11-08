package main

func (mm *MarketMaker) listenProvider() {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.stop:
			return

		case tick := <-mm.provider.Tickers():
			mm.log.Trace().Str("ask", tick.Ask.String()).Str("bid", tick.Bid.String()).Str("symbol", tick.Symbol).Msg("quote provider's tick")

			if err := mm.sendLimitsByTicker(tick); err != nil {
				mm.log.Err(err).Msg("sendLimitsByTicker")
				continue
			}

		}
	}
}
