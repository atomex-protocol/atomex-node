package synthetic

import (
	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/pkg/errors"
)

// Synthetic -
type Synthetic interface {
	Type() Type
	Ticker(tick exchange.Ticker, tickers map[string]exchange.Ticker, toSymbols map[string]string) (exchange.Ticker, error)
}

// Type -
type Type string

// types
const (
	DirectType  Type = "direct"
	DividedType Type = "divided"
)

// Config -
type Config struct {
	Symbols []string `yaml:"symbols" validate:"required"`
	Type    Type     `yaml:"type" validate:"required"`
}

func New(name string, cfg Config) (Synthetic, error) {
	switch cfg.Type {
	case DirectType:
		return NewDirect(name, cfg.Symbols...)
	case DividedType:
		return NewDivided(name, cfg.Symbols...)
	default:
		return nil, errors.Wrap(ErrUnknwonSyntheticType, string(cfg.Type))
	}
}

// errors
var (
	ErrInvalidSymbol        = errors.New("invalid tick symbol")
	ErrUnknownTicker        = errors.New("unknown ticker")
	ErrUnknwonSyntheticType = errors.New("unknown synthetic type")
)
