package strategy

import (
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy/indicators"
	"github.com/pkg/errors"
)

// Volatility -
type Volatility struct {
	std *indicators.StandardDeviation
}

// NewVolatility -
func NewVolatility(window int) *Volatility {
	return &Volatility{
		std: indicators.NewStandardDeviation(window),
	}
}

// Quotes -
func (s *Volatility) Quotes(args *Args) ([]Quote, error) {
	if !args.close.IsPositive() {
		return nil, errors.Wrapf(ErrInvalidArg, "close=%v", args.close)
	}

	s.std.Add(args.close)
	if !s.std.Full() {
		return []Quote{}, nil
	}

	// mean := s.std.Mean()
	// std := s.std.Value()

	return []Quote{}, ErrNotImplemented
}

// Is -
func (s *Volatility) Is(kind Kind) bool {
	return KindVolatility == kind
}

// Kind -
func (s *Volatility) Kind() Kind {
	return KindVolatility
}
