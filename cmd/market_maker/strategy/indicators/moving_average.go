package indicators

import "github.com/shopspring/decimal"

// MovingAverage -
type MovingAverage struct {
	values []decimal.Decimal
	min    decimal.Decimal
	max    decimal.Decimal

	window int
}

// NewMovingAverage -
func NewMovingAverage(window int) *MovingAverage {
	return &MovingAverage{
		window: window,
		values: make([]decimal.Decimal, 0, window),
		min:    decimal.Zero,
		max:    decimal.Zero,
	}
}

// Full -
func (ma *MovingAverage) Full() bool {
	return len(ma.values) == cap(ma.values)
}

// Add -
func (ma *MovingAverage) Add(value decimal.Decimal) {
	if len(ma.values) == 0 {
		ma.min = value
		ma.max = value
	} else if value.LessThan(ma.min) {
		ma.min = value
	} else if value.GreaterThan(ma.max) {
		ma.max = value
	}

	if ma.Full() {
		ma.values = append(ma.values[1:], value)
	} else {
		ma.values = append(ma.values, value)
	}
}

// Max -
func (ma *MovingAverage) Max() decimal.Decimal {
	return ma.max
}

// Min -
func (ma *MovingAverage) Min() decimal.Decimal {
	return ma.min
}

// Value -
func (ma *MovingAverage) Value() decimal.Decimal {
	sum := decimal.Zero
	for i := range ma.values {
		sum = sum.Add(ma.values[i])
	}
	count := decimal.NewFromInt(int64(len(ma.values)))
	return sum.Div(count)
}
