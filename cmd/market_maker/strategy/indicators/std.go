package indicators

import (
	"math"

	"github.com/shopspring/decimal"
)

// StandardDeviation -
type StandardDeviation struct {
	values []decimal.Decimal
	window int32
}

// NewStandardDeviation -
func NewStandardDeviation(window int) *StandardDeviation {
	return &StandardDeviation{
		values: make([]decimal.Decimal, 0, window),
		window: int32(window),
	}
}

// Full -
func (sd *StandardDeviation) Full() bool {
	return len(sd.values) == cap(sd.values)
}

// Add -
func (sd *StandardDeviation) Add(value decimal.Decimal) {
	if sd.Full() {
		sd.values = append(sd.values[1:], value)
	} else {
		sd.values = append(sd.values, value)
	}
}

// Value -
func (sd *StandardDeviation) Value() decimal.Decimal {
	mean := sd.Mean()
	value := decimal.Zero
	for i := range sd.values {
		diff := sd.values[i].Sub(mean)
		value = value.Add(diff.Pow(decimal.NewFromInt(2)))
	}

	count := decimal.NewFromInt32(int32(len(sd.values)))
	variance := value.Div(count)
	fVariance, _ := variance.Float64()
	return decimal.NewFromFloat(math.Sqrt(fVariance))
}

// Mean -
func (sd *StandardDeviation) Mean() decimal.Decimal {
	sum := decimal.Zero
	for i := range sd.values {
		sum = sum.Add(sd.values[i])
	}
	count := decimal.NewFromInt32(int32(len(sd.values)))
	return sum.Div(count)
}
