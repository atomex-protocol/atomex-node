package binance

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Base urls
const (
	BaseURLServer1 = "https://api1.binance.com"
	BaseURLServer2 = "https://api2.binance.com"
	BaseURLServer3 = "https://api3.binance.com"
)

// Error -
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`
}

// Error -
func (e Error) Error() string {
	return fmt.Sprintf("binance error: %s (%d)", e.Message, e.Code)
}

// Handle -
func (e Error) Handle() error {
	switch e.Code {
	case ErrorCodeTooManyRequests:
		time.Sleep(time.Second * 5)
		return e
	default:
		return e
	}
}

// error codes
const (
	ErrorCodeUnknown                 = -1000
	ErrorCodeDisconnected            = -1001
	ErrorCodeUnauthorized            = -1002
	ErrorCodeTooManyRequests         = -1003
	ErrorCodeServerBusy              = -1004
	ErrorCodeUnexpectedResponse      = -1006
	ErrorCodeTimeout                 = -1007
	ErrorCodeUnknownOrderComposition = -1014
	ErrorCodeTooManyOrders           = -1015
	ErrorCodeServiceShuttingDown     = -1016
	ErrorCodeUnsupportedOperation    = -1020
	ErrorCodeInvalidTimestamp        = -1021
	ErrorCodeInvalidSignature        = -1022
	ErrorCodeNotFound                = -1099

	ErrorCodeIllegalChars        = -1100
	ErrorCodeTooManyParameters   = -1101
	ErrorCodeMandatoryParamEmpty = -1102
	ErrorCodeUnknownParam        = -1103
	ErrorCodeUnreadParam         = -1104
	ErrorCodeParamEmpty          = -1105
	ErrorCodeParamNotRequired    = -1106
	ErrorCodeBadPrecision        = -1111
	ErrorCodeNoDepth             = -1112
	ErrorCodeTifNotRequired      = -1114
	ErrorCodeInvalidTif          = -1115

	ErrorCodeNewOrderRejected = -2010

	ErrorCodeExceededMaxBorrowable = -3006
)

// Server -
type Server struct {
	Time int64 `json:"serverTime"`
}

// Info -
type Info struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []Symbol      `json:"symbols"`
}

// RateLimit -
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// Symbol -
type Symbol struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuotePrecision             int      `json:"quotePrecision"`
	QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	OrderTypes                 []string `json:"orderTypes"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	Filters                    []Filter `json:"filters"`
	Permissions                []string `json:"permissions"`
}

// Filter -
type Filter struct {
	FilterType       string          `json:"filterType"`
	MinPrice         decimal.Decimal `json:"minPrice,omitempty"`
	MaxPrice         decimal.Decimal `json:"maxPrice,omitempty"`
	TickSize         decimal.Decimal `json:"tickSize,omitempty"`
	MultiplierUp     string          `json:"multiplierUp,omitempty"`
	MultiplierDown   string          `json:"multiplierDown,omitempty"`
	AvgPriceMins     int             `json:"avgPriceMins,omitempty"`
	MinQty           decimal.Decimal `json:"minQty,omitempty"`
	MaxQty           decimal.Decimal `json:"maxQty,omitempty"`
	StepSize         decimal.Decimal `json:"stepSize,omitempty"`
	MinNotional      decimal.Decimal `json:"minNotional,omitempty"`
	ApplyToMarket    bool            `json:"applyToMarket,omitempty"`
	Limit            int             `json:"limit,omitempty"`
	MaxNumOrders     int             `json:"maxNumOrders,omitempty"`
	MaxNumAlgoOrders int             `json:"maxNumAlgoOrders,omitempty"`
}

// Ticker -
type Ticker struct {
	Symbol             string          `json:"symbol"`
	PriceChange        string          `json:"priceChange"`
	PriceChangePercent string          `json:"priceChangePercent"`
	WeightedAvgPrice   string          `json:"weightedAvgPrice"`
	PrevClosePrice     string          `json:"prevClosePrice"`
	LastPrice          string          `json:"lastPrice"`
	LastQty            string          `json:"lastQty"`
	BidPrice           decimal.Decimal `json:"bidPrice"`
	AskPrice           decimal.Decimal `json:"askPrice"`
	OpenPrice          string          `json:"openPrice"`
	HighPrice          string          `json:"highPrice"`
	LowPrice           string          `json:"lowPrice"`
	Volume             string          `json:"volume"`
	QuoteVolume        string          `json:"quoteVolume"`
	OpenTime           int64           `json:"openTime"`
	CloseTime          int64           `json:"closeTime"`
	FirstID            int             `json:"firstId"`
	LastID             int             `json:"lastId"`
	Count              int             `json:"count"`
}

// Interval -
type Interval string

// intervals
const (
	IntervalMinute1  = "1m"
	IntervalMinute3  = "3m"
	IntervalMinute5  = "5m"
	IntervalMinute15 = "15m"
	IntervalMinute30 = "30m"
	IntervalHour1    = "1h"
	IntervalHour2    = "2h"
	IntervalHour4    = "4h"
	IntervalHour6    = "6h"
	IntervalHour8    = "8h"
	IntervalHour12   = "12h"
	IntervalDay1     = "1d"
	IntervalDay3     = "3d"
	IntervalWeek     = "1w"
	IntervalMonth    = "1M"
)

// OHLC -
type OHLC struct {
	OpenTime            int64
	Open                decimal.Decimal
	High                decimal.Decimal
	Low                 decimal.Decimal
	Close               decimal.Decimal
	Volume              decimal.Decimal
	CloseTime           int64
	QuoteAssetVolume    decimal.Decimal
	NumberOfTrades      uint64
	TakerBuyBaseVolume  decimal.Decimal
	TakerBuyQuoteVolume decimal.Decimal
	Ignore              decimal.Decimal
}

// UnmarshalJSON -
func (ohlc *OHLC) UnmarshalJSON(data []byte) error {
	response := []interface{}{
		&ohlc.OpenTime,
		&ohlc.Open,
		&ohlc.High,
		&ohlc.Low,
		&ohlc.Close,
		&ohlc.Volume,
		&ohlc.CloseTime,
		&ohlc.QuoteAssetVolume,
		&ohlc.NumberOfTrades,
		&ohlc.TakerBuyBaseVolume,
		&ohlc.TakerBuyQuoteVolume,
		&ohlc.Ignore,
	}
	return json.Unmarshal(data, &response)
}
