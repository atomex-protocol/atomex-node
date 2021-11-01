package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

// Exchange -
type Exchange interface {
	Start(symbols ...string) error
	Close() error
	OHLC(symbol string) ([]OHLC, error)
	Tickers() <-chan Ticker
}

// OHLC -
type OHLC struct {
	Time   time.Time
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}

// Ticker -
type Ticker struct {
	Symbol    string
	Ask       decimal.Decimal
	AskVolume decimal.Decimal
	Bid       decimal.Decimal
	BidVolume decimal.Decimal
}
