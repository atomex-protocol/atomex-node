package strategy

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Follow -
type Follow struct {
	spread Spread
}

// NewFollow -
func NewFollow(spread Spread) *Follow {
	return &Follow{spread}
}

// Quotes -
func (s *Follow) Quotes(args *Args) ([]Quote, error) {
	if !args.bid.IsPositive() {
		return nil, errors.Wrapf(ErrInvalidArg, "bid=%v", args.bid)
	}
	if !args.ask.IsPositive() {
		return nil, errors.Wrapf(ErrInvalidArg, "ask=%v", args.ask)
	}
	return []Quote{
		{
			Side:     Bid,
			Price:    args.bid.Sub(args.bid.Mul(s.spread.Bid)),
			Volume:   decimal.Zero,
			Symbol:   args.symbol,
			Strategy: KindFollow,
		},
		{
			Side:     Ask,
			Price:    args.ask.Add(args.ask.Mul(s.spread.Ask)),
			Volume:   decimal.Zero,
			Symbol:   args.symbol,
			Strategy: KindFollow,
		},
	}, nil
}

// Is -
func (s *Follow) Is(kind Kind) bool {
	return KindFollow == kind
}

// Kind -
func (s *Follow) Kind() Kind {
	return KindFollow
}
