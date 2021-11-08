package strategy

import "github.com/shopspring/decimal"

// Args -
type Args struct {
	ask       decimal.Decimal
	askVolume decimal.Decimal
	bid       decimal.Decimal
	bidVolume decimal.Decimal
	close     decimal.Decimal

	symbol string
}

// NewArgs -
func NewArgs() *Args {
	return new(Args)
}

// Bid -
func (a *Args) Bid(bid decimal.Decimal) *Args {
	a.bid = bid
	return a
}

// BidVolume -
func (a *Args) BidVolume(volume decimal.Decimal) *Args {
	a.bidVolume = volume
	return a
}

// Ask -
func (a *Args) Ask(ask decimal.Decimal) *Args {
	a.ask = ask
	return a
}

// AskVolume -
func (a *Args) AskVolume(volume decimal.Decimal) *Args {
	a.askVolume = volume
	return a
}

// Close -
func (a *Args) Close(value decimal.Decimal) *Args {
	a.close = value
	return a
}

// Symbol -
func (a *Args) Symbol(symbol string) *Args {
	a.symbol = symbol
	return a
}
