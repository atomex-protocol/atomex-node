package binance

import (
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Binance -
type Binance struct {
	ws      *Websocket
	api     *Rest
	wg      sync.WaitGroup
	stop    chan struct{}
	tickers chan exchange.Ticker

	log zerolog.Logger
}

// NewBinance -
func NewBinance(opts ...BinanceOption) *Binance {
	options := newOptions()
	for i := range opts {
		opts[i](&options)
	}

	binanceLogger := logger.New(logger.WithLogLevel(options.Level), logger.WithModuleName("binance"))

	return &Binance{
		ws:      newWebsocket(options.BaseURLWs, binanceLogger),
		api:     newRest(options.BaseURLRest, binanceLogger),
		stop:    make(chan struct{}, 1),
		tickers: make(chan exchange.Ticker, 1024),

		log: binanceLogger,
	}
}

// Start -
func (b *Binance) Start(symbols ...string) error {
	if err := b.api.Init(); err != nil {
		return errors.Wrap(err, "Binance.Init")
	}
	if err := b.ws.Connect(); err != nil {
		return errors.Wrap(err, "Binance.Connect")
	}
	if err := b.ws.SubscribeOnTickers(symbols...); err != nil {
		return errors.Wrap(err, "Binance.SubscribeOnTickers")
	}

	b.wg.Add(1)
	go b.listen()

	return nil
}

// Close -
func (b *Binance) Close() error {
	b.stop <- struct{}{}
	b.wg.Wait()

	if err := b.ws.Close(); err != nil {
		return err
	}

	close(b.tickers)
	close(b.stop)
	return nil
}

// Tickers
func (b *Binance) Tickers() <-chan exchange.Ticker {
	return b.tickers
}

// OHLC -
func (b *Binance) OHLC(symbol string) ([]exchange.OHLC, error) {
	data, err := b.api.OHLC(symbol, IntervalMinute15, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	ohlc := make([]exchange.OHLC, len(data))
	for i := range data {
		ohlc[i] = exchange.OHLC{
			Time:   time.Unix(data[i].OpenTime/1000, 0).UTC(),
			Open:   data[i].Open,
			High:   data[i].High,
			Low:    data[i].Low,
			Close:  data[i].Close,
			Volume: data[i].Volume,
		}
	}
	return ohlc, nil
}

func (b *Binance) listen() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stop:
			return
		case event := <-b.ws.Listen():
			switch typ := event.Body.(type) {
			case WebsocketTicker:
				b.tickers <- exchange.Ticker{
					Symbol:    typ.Symbol,
					Ask:       typ.Ask,
					AskVolume: typ.AskQuantity,
					Bid:       typ.Bid,
					BidVolume: typ.BidQuantity,
				}
			default:
				continue
			}
		}
	}
}
