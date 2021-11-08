package atomex

// BaseURLRestAPI -
const (
	BaseURLRestAPI       = "https://api.atomex.me"
	BaseURLRestAPIv1     = "https://api.atomex.me/v1"
	BaseURLRestAPIws     = "https://api.atomex.me/ws"
	BaseTestURLRestAPI   = "https://api.test.atomex.me"
	BaseTestURLRestAPIv1 = "https://api.test.atomex.me/v1"
	BaseTestURLRestAPIws = "https://ws.api.test.atomex.me/ws"
	WebsocketTestAPI     = "wss://ws.api.test.atomex.me/ws"
	WebsocketAPI         = "wss://ws.api.atomex.me/ws"
)

const (
	signMessage = "signing in "
)

// Sort -
type Sort string

const (
	SortDesc Sort = "Desc"
	SortAsc  Sort = "Asc"
)

// Side -
type Side string

const (
	SideBuy  Side = "Buy"
	SideSell Side = "Sell"
)

// OrderStatus -
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "Pending"
	OrderStatusPlaced          OrderStatus = "Placed"
	OrderStatusPartiallyFilled OrderStatus = "PartiallyFilled"
	OrderStatusFilled          OrderStatus = "Filled"
	OrderStatusCanceled        OrderStatus = "Canceled"
	OrderStatusRejected        OrderStatus = "Rejected"
)

// OrderType -
type OrderType string

const (
	OrderTypeReturn            OrderType = "Return"
	OrderTypeFillOrKill        OrderType = "FillOrKill"
	OrderTypeImmediateOrCancel OrderType = "ImmediateOrCancel"
	OrderTypeSolidFillOrKill   OrderType = "SolidFillOrKill"
)

// SwapStatus -
type SwapStatus string

const (
	SwapStatusCreated            SwapStatus = "Created"
	SwapStatusInvolved           SwapStatus = "Involved"
	SwapStatusPartiallyInitiated SwapStatus = "PartiallyInitiated"
	SwapStatusInitiated          SwapStatus = "Initiated"
	SwapStatusRedeemed           SwapStatus = "Redeemed"
	SwapStatusRefunded           SwapStatus = "Refunded"
	SwapStatusLost               SwapStatus = "Lost"
	SwapStatusJackpot            SwapStatus = "Jackpot"
)

// TransactionStatus -
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "Pending"
	TransactionStatusConfirmed TransactionStatus = "Confirmed"
	TransactionStatusCanceled  TransactionStatus = "Canceled"
)

// TransactionType -
type TransactionType string

const (
	TransactionTypeLock           TransactionType = "Lock"
	TransactionTypeAdditionalLock TransactionType = "AdditionalLock"
	TransactionTypeRedeem         TransactionType = "Redeem"
	TransactionTypeRefund         TransactionType = "Refund"
)

// WebsocketMethod -
type WebsocketMethod string

// methods
const (
	WebsocketMethodPing          WebsocketMethod = "ping"
	WebsocketMethodPong          WebsocketMethod = "pong"
	WebsocketMethodSubscribe     WebsocketMethod = "subscribe"
	WebsocketMethodUnsubscribe   WebsocketMethod = "unsubscribe"
	WebsocketMethodGetTopOfBook  WebsocketMethod = "getTopOfBook"
	WebsocketMethodGetSnapshot   WebsocketMethod = "getSnapshot"
	WebsocketMethodOrderSend     WebsocketMethod = "orderSend"
	WebsocketMethodOrderCancel   WebsocketMethod = "orderCancel"
	WebsocketMethodGetOrder      WebsocketMethod = "getOrder"
	WebsocketMethodGetOrders     WebsocketMethod = "getOrders"
	WebsocketMethodGetSwap       WebsocketMethod = "getSwap"
	WebsocketMethodGetSwaps      WebsocketMethod = "getSwaps"
	WebsocketMethodAddRequisites WebsocketMethod = "addRequisites"

	WebsocketMethodErrorReply         WebsocketMethod = "error"
	WebsocketMethodEntriesReply       WebsocketMethod = "entries"
	WebsocketMethodSnapshotReply      WebsocketMethod = "snapshot"
	WebsocketMethodOrderBookReply     WebsocketMethod = "orderBook"
	WebsocketMethodTopOfBookReply     WebsocketMethod = "topOfBook"
	WebsocketMethodOrderReply         WebsocketMethod = "order"
	WebsocketMethodSwapReply          WebsocketMethod = "swap"
	WebsocketMethodOrderSendReply     WebsocketMethod = "orderSendReply"
	WebsocketMethodOrderCancelReply   WebsocketMethod = "orderCancelReply"
	WebsocketMethodGetOrderReply      WebsocketMethod = "getOrderReply"
	WebsocketMethodGetOrdersReply     WebsocketMethod = "getOrdersReply"
	WebsocketMethodGetSwapReply       WebsocketMethod = "getSwapReply"
	WebsocketMethodGetSwapsReply      WebsocketMethod = "getSwapsReply"
	WebsocketMethodAddRequisitesReply WebsocketMethod = "addRequisitesReply"
)

// streams
const (
	StreamTopOfBook = "topOfBook"
	StreamOrderBook = "orderBook"
)

// WebsocketType -
type WebsocketType int

const (
	WebsocketTypeMarketData WebsocketType = iota + 1
	WebsocketTypeExchange
)

// String -
func (typ WebsocketType) String() string {
	switch typ {
	case WebsocketTypeMarketData:
		return "WebsocketTypeMarketData"
	case WebsocketTypeExchange:
		return "WebsocketTypeExchange"
	}
	return ""
}
