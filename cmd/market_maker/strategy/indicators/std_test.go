package indicators

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestStandardDeviation_Value(t *testing.T) {
	tests := []struct {
		name   string
		values []decimal.Decimal
		window int32
		want   decimal.Decimal
	}{
		{
			name: "test 1",
			values: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(4),
				decimal.NewFromInt(4),
				decimal.NewFromInt(5),
				decimal.NewFromInt(5),
				decimal.NewFromInt(7),
				decimal.NewFromInt(9),
			},
			window: 8,
			want:   decimal.NewFromInt(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StandardDeviation{
				values: tt.values,
				window: tt.window,
			}
			if got := sd.Value(); !got.Equal(tt.want) {
				t.Errorf("StandardDeviation.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
