package indicators

import "github.com/shopspring/decimal"

// ExpMovingAverage -
type ExpMovingAverage struct {
	value   decimal.Decimal
	k       decimal.Decimal
	window  int32
	counter int32
}

// NewExpMovingAverage -
func NewExpMovingAverage(window int) *ExpMovingAverage {
	return &ExpMovingAverage{
		k:      decimal.NewFromInt(2).Div(decimal.NewFromInt32(int32(window) + 1)),
		window: int32(window),
		value:  decimal.Zero,
	}
}

// Full -
func (ma *ExpMovingAverage) Full() bool {
	return ma.counter == ma.window
}

// Add -
func (ma *ExpMovingAverage) Add(value decimal.Decimal) {
	if !ma.Full() {
		ma.counter += 1
	}
	exp := decimal.NewFromInt(1).Sub(ma.k)
	ma.value = value.Mul(ma.k).Add(ma.value.Mul(exp))
}

// Value -
func (ma *ExpMovingAverage) Value() decimal.Decimal {
	return ma.value
}
