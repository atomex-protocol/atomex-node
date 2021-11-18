package main

import (
	"testing"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/stretchr/testify/assert"
)

func Test_clientOrderID_parse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    clientOrderID
		wantErr bool
	}{
		{
			name: "test 1",
			str:  "201637245985956790529XTZ_ETH",
			want: clientOrderID{
				symbol: "XTZ_ETH",
				index:  1637245985956790529,
				side:   strategy.Bid,
				kind:   strategy.KindOneByOne,
			},
		}, {
			name:    "test 2",
			str:     "201637245",
			wantErr: true,
		}, {
			name:    "test 3",
			str:     "a01637245985956790529XTZ_ETH",
			wantErr: true,
		}, {
			name: "test 4",
			str:  "2a1637245985956790529XTZ_ETH",
			want: clientOrderID{
				kind: strategy.KindOneByOne,
			},
			wantErr: true,
		}, {
			name: "test 5",
			str:  "20163724d956790529XTZ_ETH",
			want: clientOrderID{
				kind: strategy.KindOneByOne,
				side: strategy.Bid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cid clientOrderID
			if err := cid.parse(tt.str); (err != nil) != tt.wantErr {
				t.Errorf("clientOrderID.parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, cid)
		})
	}
}
