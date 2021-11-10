package strategy

import (
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy/indicators"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Volatility -
type Volatility struct {
	askStd *indicators.StandardDeviation
	bidStd *indicators.StandardDeviation

	minDist decimal.Decimal
	maxDist decimal.Decimal
	volume  decimal.Decimal

	askSpread decimal.Decimal
	bidSpread decimal.Decimal

	bid decimal.Decimal
	ask decimal.Decimal

	width decimal.Decimal

	symbol string
}

// NewVolatility -
func NewVolatility(cfg Config) *Volatility {
	return &Volatility{
		askStd: indicators.NewStandardDeviation(cfg.Window),
		bidStd: indicators.NewStandardDeviation(cfg.Window),

		volume:  cfg.Volume,
		minDist: cfg.Dist.Min,
		maxDist: cfg.Dist.Max,
		width:   cfg.Width,

		askSpread: cfg.Spread.Ask,
		bidSpread: cfg.Spread.Bid,

		bid: decimal.Zero,
		ask: decimal.Zero,

		symbol: cfg.SymbolName,
	}
}

// Quotes -
func (s *Volatility) Quotes(args *Args) ([]Quote, error) {
	if args == nil {
		return nil, errors.Wrapf(ErrInvalidArg, "nil")
	}
	if args.symbol != s.symbol {
		return nil, nil
	}
	if !args.ask.IsPositive() {
		return nil, errors.Wrapf(ErrInvalidArg, "ask=%v", args.ask)
	}
	if !args.bid.IsPositive() {
		return nil, errors.Wrapf(ErrInvalidArg, "ask=%v", args.bid)
	}

	s.askStd.Add(args.ask)
	s.bidStd.Add(args.bid)

	if !s.askStd.Full() || !s.bidStd.Full() {
		return []Quote{}, nil
	}

	s.bid = s.getBid(args.bid)
	s.ask = s.getAsk(args.ask)

	return []Quote{
		{
			Side:     Bid,
			Price:    s.bid,
			Volume:   s.volume,
			Symbol:   s.symbol,
			Strategy: KindVolatility,
		},
		{
			Side:     Ask,
			Price:    s.ask,
			Volume:   s.volume,
			Symbol:   s.symbol,
			Strategy: KindVolatility,
		},
	}, nil
}

// Is -
func (s *Volatility) Is(kind Kind) bool {
	return KindVolatility == kind
}

// Kind -
func (s *Volatility) Kind() Kind {
	return KindVolatility
}

func (s *Volatility) getAsk(ask decimal.Decimal) decimal.Decimal {
	spreadPrice := ask.Mul(decimal.NewFromInt(1).Add(s.askSpread))
	top, bottom := s.getAskChannel(ask, spreadPrice)

	if s.ask.IsZero() || ask.GreaterThan(s.ask) || bottom.GreaterThan(s.ask) || top.LessThan(s.ask) {
		return decimal.Avg(bottom, top)
	}

	return s.ask
}

func (s *Volatility) getAskChannel(ask, spreadPrice decimal.Decimal) (top decimal.Decimal, bottom decimal.Decimal) {
	askStd := s.askStd.Value()
	width := askStd.Mul(s.width)
	topWorm := ask.Add(width)
	minTop := topWorm.Add(topWorm.Mul(s.minDist))
	bottom = decimal.Max(minTop, spreadPrice)
	top = bottom.Add(bottom.Mul(s.maxDist))
	return
}

func (s *Volatility) getBid(bid decimal.Decimal) decimal.Decimal {
	spreadPrice := bid.Mul(decimal.NewFromInt(1).Sub(s.bidSpread))
	top, bottom := s.getBidChannel(bid, spreadPrice)

	if s.bid.IsZero() || bid.LessThan(s.bid) || bottom.GreaterThan(s.bid) || top.LessThan(s.bid) {
		return decimal.Avg(top, bottom)
	}

	return s.bid
}

func (s *Volatility) getBidChannel(bid, spreadPrice decimal.Decimal) (top decimal.Decimal, bottom decimal.Decimal) {
	bidStd := s.bidStd.Value()
	width := bidStd.Mul(s.width)
	bottomWorm := bid.Sub(width)
	maxBottom := bottomWorm.Sub(bottomWorm.Mul(s.minDist))
	top = decimal.Min(maxBottom, spreadPrice)
	bottom = top.Sub(top.Mul(s.maxDist))
	return
}
