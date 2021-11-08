package strategy

import "github.com/shopspring/decimal"

//  OneByOne -
type OneByOne struct {
	spread Spread
	volume decimal.Decimal
	symbol string
}

// NewOneByOne -
func NewOneByOne(cfg Config) *OneByOne {
	return &OneByOne{cfg.Spread, cfg.Volume, cfg.SymbolName}
}

// Quotes -
func (s *OneByOne) Quotes(args *Args) ([]Quote, error) {
	return []Quote{
		{
			Side:     Bid,
			Price:    decimal.NewFromInt(1).Sub(s.spread.Bid),
			Volume:   s.volume,
			Symbol:   s.symbol,
			Strategy: KindOneByOne,
		},
		{
			Side:     Ask,
			Price:    decimal.NewFromInt(1).Add(s.spread.Ask),
			Volume:   s.volume,
			Symbol:   s.symbol,
			Strategy: KindOneByOne,
		},
	}, nil
}

// Is -
func (s *OneByOne) Is(kind Kind) bool {
	return KindOneByOne == kind
}

// Kind -
func (s *OneByOne) Kind() Kind {
	return KindOneByOne
}
