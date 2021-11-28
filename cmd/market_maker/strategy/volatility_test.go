package strategy

import (
	"testing"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy/indicators"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestVolatility_Quotes(t *testing.T) {
	type fields struct {
		askStd    *indicators.StandardDeviation
		bidStd    *indicators.StandardDeviation
		minDist   decimal.Decimal
		maxDist   decimal.Decimal
		volume    decimal.Decimal
		askSpread decimal.Decimal
		bidSpread decimal.Decimal
		bid       decimal.Decimal
		ask       decimal.Decimal
		width     decimal.Decimal
		symbol    string
	}
	tests := []struct {
		name      string
		fields    fields
		askValues []decimal.Decimal
		bidValues []decimal.Decimal
		args      *Args
		want      []Quote
		wantErr   bool
	}{
		{
			name: "test 1",
			fields: fields{
				askStd:    indicators.NewStandardDeviation(4),
				bidStd:    indicators.NewStandardDeviation(4),
				minDist:   decimal.RequireFromString("0.1"),
				maxDist:   decimal.RequireFromString("0.3"),
				volume:    decimal.RequireFromString("1"),
				askSpread: decimal.RequireFromString("0.01"),
				bidSpread: decimal.RequireFromString("0.02"),
				bid:       decimal.Zero,
				ask:       decimal.Zero,
				width:     decimal.RequireFromString("1"),
				symbol:    "XTZ/ETH",
			},
			askValues: []decimal.Decimal{
				decimal.RequireFromString("2"),
				decimal.RequireFromString("3"),
				decimal.RequireFromString("3"), // std: 0.5
			},
			bidValues: []decimal.Decimal{
				decimal.RequireFromString("1"),
				decimal.RequireFromString("2"),
				decimal.RequireFromString("2"), // std: 0.5
			},
			args: new(Args).
				Ask(decimal.NewFromInt(2)).
				Bid(decimal.NewFromInt(1)).
				AskVolume(decimal.NewFromInt(2)).
				BidVolume(decimal.NewFromInt(3)).
				Symbol("XTZ/ETH"),
			want: []Quote{
				{
					Symbol:   "XTZ/ETH",
					Price:    decimal.RequireFromString("0.3825"),
					Volume:   decimal.RequireFromString("1"),
					Side:     Bid,
					Strategy: KindVolatility,
				}, {
					Symbol:   "XTZ/ETH",
					Price:    decimal.RequireFromString("3.1625"),
					Volume:   decimal.RequireFromString("1"),
					Side:     Ask,
					Strategy: KindVolatility,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Volatility{
				askStd:    tt.fields.askStd,
				bidStd:    tt.fields.bidStd,
				minDist:   tt.fields.minDist,
				maxDist:   tt.fields.maxDist,
				volume:    tt.fields.volume,
				askSpread: tt.fields.askSpread,
				bidSpread: tt.fields.bidSpread,
				bid:       tt.fields.bid,
				ask:       tt.fields.ask,
				width:     tt.fields.width,
				symbol:    tt.fields.symbol,
			}
			for i := range tt.askValues {
				s.askStd.Add(tt.askValues[i])
			}
			for i := range tt.bidValues {
				s.bidStd.Add(tt.bidValues[i])
			}
			got, err := s.Quotes(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Volatility.Quotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Len(t, got, len(tt.want)) {
				return
			}

			if !assert.Equal(t, tt.want[0].Price.String(), got[0].Price.String()) {
				return
			}

			if !assert.Equal(t, tt.want[1].Price.String(), got[1].Price.String()) {
				return
			}
		})
	}
}
