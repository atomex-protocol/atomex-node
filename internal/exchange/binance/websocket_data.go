package binance

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// urls
const (
	BaseURLWebsocket = "wss://stream.binance.com:9443/ws"
)

// WebsocketRequest -
type WebsocketRequest struct {
	Method string        `json:"method"`
	ID     uint64        `json:"id"`
	Params []interface{} `json:"params"`
}

// WebsocketError -
type WebsocketError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	ID      uint64 `json:"id"`
}

// Error -
func (we WebsocketError) Error() string {
	return fmt.Sprintf("%s (%d) with id %d", we.Message, we.Code, we.ID)
}

// methods
const (
	WebsocketMethodSubscribe         = "SUBSCRIBE"
	WebsocketMethodUnsubscribe       = "UNSUBSCRIBE"
	WebsocketMethodListSubscriptions = "LIST_SUBSCRIPTIONS"
	WebsocketMethodSetProperty       = "SET_PROPERTY"
	WebsocketMethodGetProperty       = "GET_PROPERTY"
)

// WebsocketResponse -
type WebsocketResponse struct {
	Result interface{} `json:"result"`
	ID     uint64      `json:"id"`
}

// WebsocketEvent -
type WebsocketEvent struct {
	Type string      `json:"e"`
	Time int64       `json:"E"`
	Body interface{} `json:"-"`
}

// WebsocketTicker -
type WebsocketTicker struct {
	Symbol               string          `json:"s"`
	PriceChange          decimal.Decimal `json:"p"`
	PriceChangePercent   decimal.Decimal `json:"P"`
	WeightedAveragePrice decimal.Decimal `json:"w"`
	LastPrice            decimal.Decimal `json:"c"`
	LastQuantity         decimal.Decimal `json:"Q"`
	Bid                  decimal.Decimal `json:"b"`
	BidQuantity          decimal.Decimal `json:"B"`
	Ask                  decimal.Decimal `json:"a"`
	AskQuantity          decimal.Decimal `json:"A"`
	Open                 decimal.Decimal `json:"o"`
	High                 decimal.Decimal `json:"h"`
	Low                  decimal.Decimal `json:"l"`
	FirstTradeID         int64           `json:"F"`
	LastTradeID          int64           `json:"L"`
	NumberOfTrades       int64           `json:"n"`
}

// OutboundAccountPosition -
type OutboundAccountPosition struct {
	LastUpdateTime int64 `json:"u"`
	Balances       []struct {
		Asset  string          `json:"a"`
		Free   decimal.Decimal `json:"f"`
		Locked decimal.Decimal `json:"l"`
	} `json:"B"`
}

// BalanceUpdate -
type BalanceUpdate struct {
	Asset        string          `json:"a"`
	BalanceDelta decimal.Decimal `json:"d"`
	ClearTime    decimal.Decimal `json:"T"`
}

// KLine -
type KLine struct {
	Symbol string `json:"s"`
	KLine  struct {
		OpenTime     int64           `json:"t"`
		CloseTime    int64           `json:"T"`
		Symbol       string          `json:"s"`
		Interval     string          `json:"i"`
		FirstTradeID int64           `json:"f"`
		LastTradeID  int64           `json:"L"`
		Open         decimal.Decimal `json:"o"`
		Close        decimal.Decimal `json:"c"`
		High         decimal.Decimal `json:"h"`
		Low          decimal.Decimal `json:"l"`
		BaseVolume   decimal.Decimal `json:"v"`
		QuoteVolume  decimal.Decimal `json:"q"`
		TradesNumber int64           `json:"n"`
		IsClosed     bool            `json:"x"`
	} `json:"k"`
}

// ExecutionReport -
type ExecutionReport struct {
	Symbol                   string          `json:"s"`
	ClientOrderID            string          `json:"c"`
	Side                     string          `json:"S"`
	OrderType                string          `json:"o"`
	TimeInForce              string          `json:"f"`
	Quantity                 decimal.Decimal `json:"q"`
	Price                    decimal.Decimal `json:"p"`
	StopPrice                decimal.Decimal `json:"P"`
	IcebergQuantity          decimal.Decimal `json:"F"`
	OrderListID              int64           `json:"g"`
	OriginalClientOrderID    string          `json:"C"`
	ExecutionType            string          `json:"x"`
	OrderStatus              string          `json:"X"`
	RejectReason             string          `json:"r"`
	OrderID                  int64           `json:"i"`
	LastExecutedQuantity     decimal.Decimal `json:"l"`
	CumulativeFilledQuantity decimal.Decimal `json:"z"`
	LastExecutedPrice        decimal.Decimal `json:"L"`
	CommisionAmount          decimal.Decimal `json:"n"`
	CommisionAsset           string          `json:"N,omitempty"`
	TransactionTime          int64           `json:"T"`
	TradeID                  int64           `json:"t"`
	OnBook                   bool            `json:"w"`
	IsMaker                  bool            `json:"m"`
	OrderCreationTime        int64           `json:"O"`
}

// ListStatus -
type ListStatus struct {
	Symbol            string `json:"s"`
	OrderListID       int64  `json:"g"`
	ContingencyType   string `json:"c"`
	ListStatusType    string `json:"l"`
	ListOrderStatus   string `json:"L"`
	ListRejectReason  string `json:"r"`
	ListClientOrderID string `json:"C"`
	TransactionTime   int64  `json:"T"`
	Objects           []struct {
		Symbol        string `json:"s"`
		OrderID       int    `json:"i"`
		ClientOrderID string `json:"c"`
	} `json:"O"`
}

// BookTicker -
type BookTicker struct {
	UpdateID  int64           `json:"u"`
	Symbol    string          `json:"s"`
	Bid       decimal.Decimal `json:"b"`
	BidVolume decimal.Decimal `json:"B"`
	Ask       decimal.Decimal `json:"a"`
	AskVolume decimal.Decimal `json:"A"`
}

// UnmarshalJSON -
func (we *WebsocketEvent) UnmarshalJSON(data []byte) error {
	type event WebsocketEvent
	var e event
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	we.Type = e.Type

	switch we.Type {
	case "24hrTicker":
		var ticker WebsocketTicker
		if err := json.Unmarshal(data, &ticker); err != nil {
			return err
		}
		we.Body = ticker
	case "kline":
		var res KLine
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	case "outboundAccountPosition":
		var res OutboundAccountPosition
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	case "balanceUpdate":
		var res BalanceUpdate
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	case "executionReport":
		var res ExecutionReport
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	case "listStatus":
		var res ListStatus
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	case "":
		var res BookTicker
		if err := json.Unmarshal(data, &res); err != nil {
			return err
		}
		we.Body = res
	default:
		return errors.Errorf("unknown event type: %s", e.Type)
	}
	return nil
}
