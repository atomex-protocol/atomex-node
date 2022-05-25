package main

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_amountToInt(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name     string
		value    string
		decimals int
		want     string
	}{
		{
			name:     "test 1",
			value:    "0.9876541",
			decimals: 6,
			want:     "987654",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := decimal.NewFromString(tt.value)
			if got := amountToInt(value, tt.decimals); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("amountToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
