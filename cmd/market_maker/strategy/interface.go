package strategy

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Strategy -
type Strategy interface {
	Quotes(args *Args) ([]Quote, error)
}

// New -
func New(cfg Config) (Strategy, error) {
	switch cfg.Kind {
	case KindFollow:
		return NewFollow(cfg.Spread), nil
	case KindOneByOne:
		return NewOneByOne(cfg), nil
	case KindVolatility:
		if cfg.Window < 1 {
			return nil, errors.Wrapf(ErrInvalidArg, "window=%d", cfg.Window)
		}
		return NewVolatility(cfg.Window), nil
	default:
		return nil, errors.Wrap(ErrUnknownStrategy, string(cfg.Kind))
	}
}

// Kind -
type Kind string

// kinds
const (
	KindFollow     = "follow"
	KindOneByOne   = "one-by-one"
	KindVolatility = "volatility"
)

// Spread -
type Spread struct {
	Ask decimal.Decimal `yaml:"ask"`
	Bid decimal.Decimal `yaml:"bid"`
}

// UnmarshalYAML -
func (s *Spread) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type buf Spread
	if err := unmarshal((*buf)(s)); err == nil {
		return nil
	}

	var offset decimal.Decimal
	if err := unmarshal(&offset); err != nil {
		return err
	}

	s.Ask = offset
	s.Bid = offset
	return nil
}

// Quote -
type Quote struct {
	Symbol string
	Side   Side
	Price  decimal.Decimal
	Volume decimal.Decimal
}

// Side -
type Side int

// sides
const (
	Bid Side = iota
	Ask
)

// Config -
type Config struct {
	SymbolName string          `yaml:"symbol"`
	Kind       Kind            `yaml:"kind"`
	Spread     Spread          `yaml:"spread"`
	Volume     decimal.Decimal `yaml:"volume"`
	Window     int             `yaml:"window"`
}
