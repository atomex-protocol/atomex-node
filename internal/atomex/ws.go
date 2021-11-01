package atomex

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Websocket -
type Websocket struct {
	url       *url.URL
	uri       string
	restURI   string
	conn      *websocket.Conn
	requestID uint64
	algo      string
	token     TokenResponse
	logger    zerolog.Logger

	errorChan chan error
	msgs      chan Message
	stop      chan struct{}
	wg        sync.WaitGroup

	// test bool
}

// Exchange -
type Exchange struct {
	*Websocket
}

// Market -
type Market struct {
	*Websocket
}

// NewMarket -
func NewMarket(opts ...WebsocketOption) (*Market, error) {
	ws, err := NewWebsocket(WebsocketTypeMarketData, opts...)
	if err != nil {
		return nil, err
	}
	return &Market{ws}, nil
}

// NewExchange -
func NewExchange(opts ...WebsocketOption) (*Exchange, error) {
	ws, err := NewWebsocket(WebsocketTypeExchange, opts...)
	if err != nil {
		return nil, err
	}
	return &Exchange{ws}, nil
}

// NewWebsocket -
func NewWebsocket(typ WebsocketType, opts ...WebsocketOption) (*Websocket, error) {
	ws := Websocket{
		algo:      signers.AlgorithmBlake2bWithEcdsaSecp256k1,
		errorChan: make(chan error, 1024),
		msgs:      make(chan Message, 1024),
		stop:      make(chan struct{}, 2),
		logger:    logger.New(logger.WithModuleName(typ.String())),
		uri:       WebsocketAPI,
		restURI:   BaseURLRestAPIws,
	}

	for i := range opts {
		opts[i](&ws)
	}

	switch typ {
	case WebsocketTypeExchange:
		ws.uri = fmt.Sprintf("%s/exchange", ws.uri)
	case WebsocketTypeMarketData:
		ws.uri = fmt.Sprintf("%s/marketdata", ws.uri)
	default:
		return nil, errors.Errorf("unknwon websocket type: %d", typ)
	}

	u, err := url.Parse(ws.uri)
	if err != nil {
		return nil, err
	}
	ws.url = u

	return &ws, nil
}

// Connect -
func (ws *Websocket) Connect(keys *signers.Key) error {
	token, err := NewRest(WithURL(ws.restURI), WithSignatureAlgorithm(ws.algo)).Token(keys)
	if err != nil {
		return errors.Wrap(err, "Token")
	}
	ws.token = token

	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Bearer %s", ws.token.Token))
	header.Add("Content-Type", "application/json")
	c, _, err := websocket.DefaultDialer.Dial(ws.url.String(), header)
	if err != nil {
		return errors.Wrap(err, "Connect Dial")
	}
	ws.conn = c

	ws.wg.Add(1)
	go ws.ping()

	ws.wg.Add(1)
	go ws.listen()

	return nil
}

// Close -
func (ws *Websocket) Close() error {
	ws.stop <- struct{}{} // for listen
	return ws.close()
}

func (ws *Websocket) close() error {
	ws.stop <- struct{}{} // for ping
	ws.wg.Wait()

	if err := ws.conn.Close(); err != nil {
		return err
	}

	close(ws.errorChan)
	close(ws.msgs)
	close(ws.stop)
	return nil
}

// Errors - channels with protocol errors
func (ws *Websocket) Errors() <-chan error {
	return ws.errorChan
}

// Listen - channels with messages
func (ws *Websocket) Listen() <-chan Message {
	return ws.msgs
}

func (ws *Websocket) ping() {
	defer ws.wg.Done()

	keepAliveTicker := time.NewTicker(time.Second * 15)
	defer keepAliveTicker.Stop()

	for {
		select {
		case <-ws.stop:
			return
		case <-keepAliveTicker.C:
			if err := ws.send(WebsocketMessage{
				Method: WebsocketMethodPing,
			}); err != nil {
				if err := ws.reconnect(); err != nil {
					ws.logger.Err(err).Msg("reconnect")
					ws.logger.Warn().Msg("retry after 5 seconds")
					time.Sleep(time.Second * 5)
				}
			}
		}
	}
}

func (ws *Websocket) listen() {
	defer ws.wg.Done()

	for {
		select {
		case <-ws.stop:
			return
		default:
			if err := ws.readAllMessages(); err != nil {
				switch {
				case errors.Is(err, ErrTimeout) || websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
					if err := ws.reconnect(); err != nil {
						ws.logger.Err(err).Msg("reconnect")
						ws.logger.Warn().Msg("retry after 5 seconds")
						time.Sleep(time.Second * 5)
					}
				case websocket.IsCloseError(err, websocket.CloseNormalClosure):
					ws.logger.Err(err).Msg("readAllMessages")
					if err := ws.close(); err != nil {
						ws.logger.Err(err).Msg("close")
					}
					return
				case errors.Is(err, ErrEmptyResponse):
				default:
					ws.logger.Err(err).Msg("readAllMessages")
				}
			}
		}
	}
}

func (ws *Websocket) reconnect() error {
	ws.logger.Warn().Msg("reconnecting...")

	if err := ws.conn.Close(); err != nil {
		return errors.Wrap(err, "reconnect Close")
	}
	ws.logger.Trace().Msg("connection closed")

	c, _, err := websocket.DefaultDialer.Dial(ws.url.String(), http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ws.token.Token)},
	})
	if err != nil {
		return errors.Wrap(err, "reconnect Dial")
	}
	ws.conn = c
	return nil
}

func (ws *Websocket) readAllMessages() error {
	if err := ws.conn.SetReadDeadline(time.Now().Add(time.Second * 20)); err != nil {
		return errors.Wrap(err, "SetReadDeadline")
	}

	_, reader, err := ws.conn.NextReader()
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return ErrTimeout
		}
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
			return err
		}

		return errors.Wrap(err, "NextReader")
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	ws.logger.Trace().RawJSON("data", data).Msg("server->client")

	var msg WebsocketResponse
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	return ws.handleMessage(msg)
}

func (ws *Websocket) handleMessage(msg WebsocketResponse) error {
	if msg.Data == nil {
		return nil
	}

	switch msg.Event {
	case WebsocketMethodEntriesReply:
		var entries []Entry
		if err := json.Unmarshal(*msg.Data, &entries); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, entries}
	case WebsocketMethodSnapshotReply:
		var snapshot Snapshot
		if err := json.Unmarshal(*msg.Data, &snapshot); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, snapshot}
	case WebsocketMethodTopOfBookReply:
		var reply []TopOfBook
		if err := json.Unmarshal(*msg.Data, &reply); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, reply}
	case WebsocketMethodOrderBookReply:
		var entries []OrderBookItemWebsocket
		if err := json.Unmarshal(*msg.Data, &entries); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, entries}
	case WebsocketMethodOrderReply:
		var order Order
		if err := json.Unmarshal(*msg.Data, &order); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, order}
	case WebsocketMethodSwapReply:
		var swap Swap
		if err := json.Unmarshal(*msg.Data, &swap); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, swap}
	case WebsocketMethodOrderSendReply:
		orderID, err := strconv.ParseInt(string(*msg.Data), 10, 64)
		if err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, orderID}
	case WebsocketMethodOrderCancelReply:
		orderID, err := strconv.ParseInt(string(*msg.Data), 10, 64)
		if err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, orderID}
	case WebsocketMethodGetOrderReply:
		var order Order
		if err := json.Unmarshal(*msg.Data, &order); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, order}
	case WebsocketMethodGetOrdersReply:
		var orders []Order
		if err := json.Unmarshal(*msg.Data, &orders); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, orders}
	case WebsocketMethodGetSwapReply:
		var swap Swap
		if err := json.Unmarshal(*msg.Data, &swap); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, swap}
	case WebsocketMethodGetSwapsReply:
		var swaps []Swap
		if err := json.Unmarshal(*msg.Data, &swaps); err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, swaps}
	case WebsocketMethodAddRequisitesReply:
		success, err := strconv.ParseBool(string(*msg.Data))
		if err != nil {
			return err
		}
		ws.msgs <- Message{msg.Event, success}
	case WebsocketMethodErrorReply:
		if msg.Data != nil {
			ws.errorChan <- errors.New(string(*msg.Data))
		}

	default:
		return errors.Errorf("unsupported event: %s", msg.Event)
	}
	return nil
}

func (ws *Websocket) send(msg WebsocketMessage) error {
	ws.requestID++
	msg.RequestID = ws.requestID

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ws.logger.Trace().RawJSON("data", data).Msg("client->server")
	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *Websocket) subscribe(stream string) error {
	return ws.send(newWebsocketMessage(WebsocketMethodSubscribe, toQuotedBytes(stream)))
}

func (ws *Websocket) unsubscribe(stream string) error {
	return ws.send(newWebsocketMessage(WebsocketMethodUnsubscribe, toQuotedBytes(stream)))
}

// SubscribeToOrderBook -
func (market *Market) SubscribeToOrderBook() error {
	return market.subscribe(StreamOrderBook)
}

// UnsubscribeFromOrderBook -
func (market *Market) UnsubscribeFromOrderBook() error {
	return market.unsubscribe(StreamOrderBook)
}

// SubscribeToTopOfBook -
func (market *Market) SubscribeToTopOfBook() error {
	return market.subscribe(StreamTopOfBook)
}

// UnsubscribeFromTopOfBook -
func (market *Market) UnsubscribeFromTopOfBook() error {
	return market.unsubscribe(StreamTopOfBook)
}

// GetTopOfBook -
func (market *Market) GetTopOfBook(symbols ...string) error {
	if len(symbols) == 0 {
		return errors.Wrap(ErrInvalidArg, "array of symbols is empty in GetTopOfBook")
	}
	data, err := json.Marshal(symbols)
	if err != nil {
		return err
	}

	return market.send(newWebsocketMessage(WebsocketMethodGetTopOfBook, data))
}

// GetSnapshot -
func (market *Market) GetSnapshot(symbol string) error {
	if len(symbol) == 0 {
		return errors.Wrap(ErrInvalidArg, "symbol is empty in GetSnapshot")
	}
	return market.send(newWebsocketMessage(WebsocketMethodGetSnapshot, toQuotedBytes(symbol)))
}

// SendOrder -
func (ex *Exchange) SendOrder(order AddOrderRequest) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodOrderSend, data))
}

// CancelOrder -
func (ex *Exchange) CancelOrder(order CancelOrderRequest) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodOrderCancel, data))
}

// Order -
func (ex *Exchange) Order(id int64) error {
	data, err := json.Marshal(GetByIDRequest{id})
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodGetOrder, data))
}

// Order -
func (ex *Exchange) Orders(req OrdersRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodGetOrders, data))
}

// Swap -
func (ex *Exchange) Swap(id int64) error {
	data, err := json.Marshal(GetByIDRequest{id})
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodGetSwap, data))
}

// Swaps -
func (ex *Exchange) Swaps(req SwapsWebsocketRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodGetSwaps, data))
}

// AddRequisites -
func (ex *Exchange) AddRequisites(req AddSwapRequisitesWebsocketRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return ex.send(newWebsocketMessage(WebsocketMethodGetSwaps, data))
}
