package synthetic

import (
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Divided -
type Divided struct {
	name   string
	first  string
	second string
}

// NewDivided -
func NewDivided(name string, symbols ...string) (*Divided, error) {
	if len(symbols) != 2 {
		return nil, errors.Errorf("invalid symbols count in divided synthetic config: %v", symbols)
	}

	return &Divided{name, symbols[0], symbols[1]}, nil
}

// Type -
func (d *Divided) Type() Type {
	return DividedType
}

// Ticker -
func (d *Divided) Ticker(tick exchange.Ticker, tickers map[string]exchange.Ticker, toSymbols map[string]string) (exchange.Ticker, error) {
	if d.first != tick.Symbol && d.second != tick.Symbol {
		return tick, errors.Wrap(ErrInvalidSymbol, tick.Symbol)
	}
	firstSymbol, ok := toSymbols[d.first]
	if !ok {
		return tick, errors.Wrap(ErrInvalidSymbol, d.first)
	}
	first, ok := tickers[firstSymbol]
	if !ok {
		return first, errors.Wrap(ErrUnknownTicker, firstSymbol)
	}

	secondSymbol, ok := toSymbols[d.second]
	if !ok {
		return tick, errors.Wrap(ErrInvalidSymbol, d.second)
	}
	second, ok := tickers[secondSymbol]
	if !ok {
		return first, errors.Wrap(ErrUnknownTicker, secondSymbol)
	}

	return exchange.Ticker{
		Symbol:    d.name,
		Ask:       first.Ask.Div(second.Bid),
		Bid:       first.Bid.Div(second.Ask),
		AskVolume: decimal.Zero, // TODO: compute volume
		BidVolume: decimal.Zero,
	}, nil
}
