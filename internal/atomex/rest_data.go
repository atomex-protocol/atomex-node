package atomex

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/shopspring/decimal"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Errors  []struct {
		Field   string `json:"field,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

// Error -
func (e Error) Error() string {
	var builder strings.Builder
	builder.WriteString(e.Message)
	builder.WriteString(" (")
	builder.WriteString(strconv.FormatInt(int64(e.Code), 10))
	builder.WriteString(")")
	for i := range e.Errors {
		builder.WriteString("\r\n")
		builder.WriteString(e.Errors[i].Field)
		builder.WriteString(": ")
		builder.WriteString(e.Errors[i].Message)
	}
	return builder.String()
}

// TokenRequest -
type TokenRequest struct {
	Timestamp int64  `json:"timeStamp"`
	Message   string `json:"message"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Algorithm string `json:"algorithm"`
}

// NewTokenRequest -
func NewTokenRequest(message, algorithm string, publicKey []byte) TokenRequest {
	return TokenRequest{
		Timestamp: time.Now().UnixNano() / 1_000_000,
		PublicKey: hex.EncodeToString(publicKey),
		Message:   message,
		Algorithm: algorithm,
	}
}

// Sign -
func (req *TokenRequest) Sign(key *signers.Key) error {
	signingMessage := fmt.Sprintf("%s%d", req.Message, req.Timestamp)
	signature, err := signers.Sign(req.Algorithm, key, []byte(signingMessage), true)
	if err != nil {
		return err
	}
	req.Signature = hex.EncodeToString(signature)

	return nil
}

// TokenResponse -
type TokenResponse struct {
	ID      string `json:"id"`
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// TopOfBook -
type TopOfBook struct {
	Symbol    string          `json:"symbol"`
	Timestamp int64           `json:"timeStamp"`
	Bid       decimal.Decimal `json:"bid"`
	Ask       decimal.Decimal `json:"ask"`
}

// OrderBook -
type OrderBook struct {
	UpdateID int64           `json:"updateId"`
	Symbol   string          `json:"symbol"`
	Entries  []OrderBookItem `json:"entries"`
}

// OrderBookItem -
type OrderBookItem struct {
	Side       Side  `json:"side"`
	Price      int   `json:"price"`
	QtyProfile []int `json:"qtyProfile"`
}

// AddOrderRequest -
type AddOrderRequest struct {
	ClientOrderID string         `json:"clientOrderId"`
	Symbol        string         `json:"symbol"`
	Price         float64        `json:"price"`
	Qty           float64        `json:"qty"`
	Side          Side           `json:"side"`
	Type          OrderType      `json:"type"`
	ProofsOfFunds []ProofOfFunds `json:"proofsOfFunds,omitempty"`
	Requisites    *Requisites    `json:"requisites,omitempty"`
}

// ProofOfFunds -
type ProofOfFunds struct {
	Address   string `json:"address"`
	Currency  string `json:"currency"`
	Timestamp int64  `json:"timeStamp"`
	Message   string `json:"message"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Algorithm string `json:"algorithm"`
}

// AddOrderResponse -
type AddOrderResponse struct {
	OrderID int64 `json:"orderId"`
}

// OrdersRequest -
type OrdersRequest struct {
	Symbols []string `json:"symbols,omitempty"`
	Sort    Sort     `json:"sort,omitempty"`
	Offset  uint64   `json:"offset,omitempty"`
	Limit   uint64   `json:"limit,omitempty"`
	Active  bool     `json:"onlyActive,omitempty"`
}

func (req OrdersRequest) getArgs() url.Values {
	args := make(url.Values)
	if len(req.Symbols) > 0 {
		args.Add("symbols", strings.Join(req.Symbols, ","))
	}
	if req.Sort != "" {
		args.Add("sort", string(req.Sort))
	}
	if req.Offset > 0 && req.Offset <= 2147483647 {
		args.Add("offset", strconv.FormatUint(req.Offset, 10))
	}
	if req.Limit > 0 && req.Offset <= 10000 {
		args.Add("limit", strconv.FormatUint(req.Limit, 10))
	}
	if req.Active {
		args.Add("active", strconv.FormatBool(req.Active))
	}

	return args
}

// Order -
type Order struct {
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
	Swaps         []Swap          `json:"swaps"`
}

// Swap -
type Swap struct {
	ID           int64           `json:"id"`
	Symbol       string          `json:"symbol"`
	Side         Side            `json:"side"`
	TimeStamp    time.Time       `json:"timeStamp"`
	Price        decimal.Decimal `json:"price"`
	Qty          decimal.Decimal `json:"qty"`
	Secret       string          `json:"secret"`
	SecretHash   string          `json:"secretHash"`
	IsInitiator  bool            `json:"isInitiator"`
	User         User            `json:"user"`
	CounterParty User            `json:"counterParty"`
}

// User -
type User struct {
	Requisites   Requisites    `json:"requisites"`
	Status       SwapStatus    `json:"status"`
	Trades       []Trade       `json:"trades"`
	Transactions []Transaction `json:"transactions"`
}

// Requisites -
type Requisites struct {
	SecretHash            string  `json:"secretHash,omitempty"`
	ReceivingAddress      string  `json:"receivingAddress,omitempty"`
	RefundAddress         string  `json:"refundAddress,omitempty"`
	RewardForRedeem       float64 `json:"rewardForRedeem,omitempty"`
	LockTime              int64   `json:"lockTime,omitempty"`
	BaseCurrencyContract  string  `json:"baseCurrencyContract"`
	QuoteCurrencyContract string  `json:"quoteCurrencyContract"`
}

// Trade -
type Trade struct {
	OrderID int64           `json:"orderId"`
	Price   decimal.Decimal `json:"price"`
	Qty     decimal.Decimal `json:"qty"`
}

// Transaction -
type Transaction struct {
	Currency      string            `json:"currency"`
	TxID          string            `json:"txId"`
	BlockHeight   int64             `json:"blockHeight"`
	Confirmations int64             `json:"confirmations"`
	Status        TransactionStatus `json:"status"`
	Type          TransactionType   `json:"type"`
}

// DefaultResponse -
type DefaultResponse struct {
	Result bool `json:"result"`
}

// SwapsRequest -
type SwapsRequest struct {
	Symbols   []string
	Sort      Sort
	Offset    uint64
	Limit     uint64
	AfterID   int64
	Active    bool
	Completed bool
}

func (req SwapsRequest) getArgs() url.Values {
	args := make(url.Values)
	if len(req.Symbols) > 0 {
		args.Add("symbols", strings.Join(req.Symbols, ","))
	}
	if req.Sort != "" {
		args.Add("sort", string(req.Sort))
	} else {
		args.Add("sort", string(SortDesc))
	}
	if req.Offset > 0 && req.Offset <= 2147483647 {
		args.Add("offset", strconv.FormatUint(req.Offset, 10))
	}
	if req.Limit > 0 && req.Offset <= 10000 {
		args.Add("limit", strconv.FormatUint(req.Limit, 10))
	}
	if req.Active {
		args.Add("active", strconv.FormatBool(req.Active))
	}
	if req.AfterID > 0 {
		args.Add("afterId", strconv.FormatInt(req.AfterID, 10))
	}
	if req.Completed {
		args.Add("completed", strconv.FormatBool(req.Completed))
	}

	return args
}

// SwapsWebsocketRequest -
type SwapsWebsocketRequest struct {
	Symbols []string `json:"symbols.omitempty"`
	Sort    Sort     `json:"sort,omitempty"`
	Offset  uint64   `json:"offset,omitempty"`
	Limit   uint64   `json:"limit,omitempty"`
}

// AddSwapRequisitesRequest -
type AddSwapRequisitesRequest struct {
	SecretHash       string  `json:"secretHash"`
	ReceivingAddress string  `json:"receivingAddress"`
	RefundAddress    string  `json:"refundAddress"`
	RewardForRedeem  float64 `json:"rewardForRedeem"`
	LockTime         int64   `json:"lockTime"`
}

// AddSwapRequisitesWebsocketRequest -
type AddSwapRequisitesWebsocketRequest struct {
	ID               int64   `json:"id"`
	SecretHash       string  `json:"secretHash"`
	ReceivingAddress string  `json:"receivingAddress"`
	RefundAddress    string  `json:"refundAddress"`
	RewardForRedeem  float64 `json:"rewardForRedeem"`
	LockTime         int64   `json:"lock_time"`
}

// SymbolInfo -
type SymbolInfo struct {
	Name       string          `json:"name"`
	MinimumQty decimal.Decimal `json:"minimumQty"`
}

// CancelOrderRequest -
type CancelOrderRequest struct {
	ID     int64  `json:"id"`
	Symbol string `json:"symbol"`
	Side   Side   `json:"side"`
}

// GetByIDRequest -
type GetByIDRequest struct {
	ID int64 `json:"id"`
}
