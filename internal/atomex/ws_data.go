package atomex

import (
	stdJSON "encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// WebsocketMessage -
type WebsocketMessage struct {
	Method    WebsocketMethod     `json:"method"`
	Data      *stdJSON.RawMessage `json:"data,omitempty"`
	RequestID uint64              `json:"requestId"`
}

func newWebsocketMessage(method WebsocketMethod, data []byte) WebsocketMessage {
	wm := WebsocketMessage{
		Method: method,
	}

	if len(data) > 0 {
		raw := stdJSON.RawMessage(data)
		wm.Data = &raw
	}

	return wm
}

// WebsocketResponse -
type WebsocketResponse struct {
	Event     WebsocketMethod     `json:"event"`
	Data      *stdJSON.RawMessage `json:"data,omitempty"`
	RequestID uint64              `json:"requestId"`
}

// MarketData -
type MarketData struct {
	UpdateID int64  `json:"updateId"`
	Symbol   string `json:"symbol"`
}

// Entry -
type Entry struct {
	*MarketData
	Side       Side              `json:"side"`
	Price      decimal.Decimal   `json:"price"`
	QtyProfile []decimal.Decimal `json:"qtyProfile"`
}

// OrderBookItemWebsocket -
type OrderBookItemWebsocket struct {
	UpdateID   int64             `json:"UpdateId"`
	Symbol     string            `json:"Symbol"`
	Side       Side              `json:"Side"`
	Price      decimal.Decimal   `json:"Price"`
	QtyProfile []decimal.Decimal `json:"QtyProfile"`
}

// Snapshot -
type Snapshot struct {
	*MarketData
	Entries []Entry `json:"entries"`
}

// TopOfBook -
type TopOfBookWebsocket struct {
	Symbol    string          `json:"Symbol"`
	Timestamp int64           `json:"TimeStamp"`
	Ask       decimal.Decimal `json:"Ask"`
	Bid       decimal.Decimal `json:"Bid"`
}

// Message -
type Message struct {
	Event WebsocketMethod
	Value interface{}
}

// AddOrderWebsocketResponse -
type AddOrderWebsocketResponse struct {
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
}

// CancelOrderWebsocketResponse -
type CancelOrderWebsocketResponse struct {
	OrderID int64 `json:"orderId"`
	Result  bool  `json:"result"`
}

// OrderWebsocket -
type OrderWebsocket struct {
	ID            int64           `json:"id"`
	ClientOrderID string          `json:"clientOrderId"`
	Symbol        string          `json:"symbol"`
	Side          Side            `json:"side"`
	Timestamp     time.Time       `json:"timeStamp"`
	Price         decimal.Decimal `json:"price"`
	Qty           decimal.Decimal `json:"qty"`
	LeaveQty      decimal.Decimal `json:"leaveQty"`
	Type          OrderType       `json:"type"`
	Status        OrderStatus     `json:"status"`
	Trades        []Trade         `json:"trades"`
	Swaps         []int64         `json:"swaps"`
}
