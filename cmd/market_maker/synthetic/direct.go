package synthetic

import (
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/pkg/errors"
)

// Direct -
type Direct struct {
	name   string
	symbol string
}

// NewDirect -
func NewDirect(name string, symbols ...string) (*Direct, error) {
	if len(symbols) != 1 {
		return nil, errors.Errorf("invalid symbols count in direct synthetic config: %v", symbols)
	}
	return &Direct{name, symbols[0]}, nil
}

// Type -
func (d *Direct) Type() Type {
	return DirectType
}

// Ticker -
func (d *Direct) Ticker(tick exchange.Ticker, tickers map[string]exchange.Ticker, toSymbols map[string]string) (exchange.Ticker, error) {
	if d.symbol != tick.Symbol {
		return tick, errors.Wrap(ErrInvalidSymbol, tick.Symbol)
	}

	return exchange.Ticker{
		Symbol:    d.name,
		Ask:       tick.Ask,
		AskVolume: tick.AskVolume,
		Bid:       tick.Bid,
		BidVolume: tick.BidVolume,
	}, nil
}
