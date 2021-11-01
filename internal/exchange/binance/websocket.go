package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Websocket -
type Websocket struct {
	url           string
	conn          *websocket.Conn
	log           zerolog.Logger
	subscriptions map[uint64]WebsocketRequest
	id            uint64

	reconnectTimeout time.Duration
	readTimeout      time.Duration

	events  chan WebsocketEvent
	connect chan struct{}
	stop    chan struct{}

	wg sync.WaitGroup
}

func newWebsocket(url string, logger zerolog.Logger) *Websocket {
	return &Websocket{
		url:              url,
		log:              logger,
		reconnectTimeout: time.Second,
		readTimeout:      15 * time.Second,
		subscriptions:    make(map[uint64]WebsocketRequest),
		connect:          make(chan struct{}, 1),
		events:           make(chan WebsocketEvent, 1024),
		stop:             make(chan struct{}, 1),
		id:               1,
	}
}

func (ws *Websocket) dial() error {
	dialer := websocket.Dialer{
		Subprotocols:    []string{"p1", "p2"},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Proxy:           http.ProxyFromEnvironment,
	}

	c, resp, err := dialer.Dial(ws.url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ws.conn = c
	return nil
}

// Connect -
func (ws *Websocket) Connect() error {
	if err := ws.dial(); err != nil {
		return err
	}

	ws.wg.Add(1)
	go ws.listen()

	return nil
}

// Listen -
func (ws *Websocket) Listen() <-chan WebsocketEvent {
	return ws.events
}

// Close -
func (ws *Websocket) Close() error {
	ws.stop <- struct{}{}
	ws.wg.Wait()

	for _, sub := range ws.subscriptions {
		if err := ws.send(&WebsocketRequest{
			Method: WebsocketMethodUnsubscribe,
			Params: sub.Params,
		}); err != nil {
			return err
		}
	}

	if ws.conn != nil {
		if err := ws.conn.Close(); err != nil {
			return err
		}
	}

	close(ws.stop)
	close(ws.events)
	close(ws.connect)
	return nil
}

func (ws *Websocket) reconnect() error {
	time.Sleep(ws.reconnectTimeout)

	ws.log.Warn().Msg("reconnecting...")

	if err := ws.dial(); err != nil {
		return errors.Wrap(err, "dial")
	}

	if err := ws.resubscribe(); err != nil {
		return errors.Wrap(err, "resubscribe")
	}
	ws.log.Warn().Msg("reconnected")
	return nil
}

func (ws *Websocket) listen() {
	defer ws.wg.Done()

	if ws.conn == nil {
		return
	}

	if err := ws.conn.SetReadDeadline(time.Now().Add(ws.readTimeout)); err != nil {
		ws.log.Err(err).Msg("SetReadDeadline")
		return
	}

	for {
		select {
		case <-ws.stop:
			return
		case <-ws.connect:
			if err := ws.reconnect(); err != nil {
				ws.log.Err(err).Msg("ReadMessage")
				ws.connect <- struct{}{}
			}
		default:
			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				ws.log.Err(err).Msg("ReadMessage")
				ws.connect <- struct{}{}
				continue
			}

			if err := ws.conn.SetReadDeadline(time.Now().Add(ws.readTimeout)); err != nil {
				ws.log.Err(err).Msg("SetReadDeadline")
				ws.connect <- struct{}{}
				continue
			}

			ws.log.Trace().RawJSON("data", msg).Msg("server->client")

			if err := ws.handleMessage(msg); err != nil {
				ws.log.Err(err).Msg("handleMessage")
			}
		}
	}
}

func (ws *Websocket) resubscribe() error {
	for _, sub := range ws.subscriptions {
		if err := ws.send(&sub); err != nil {
			return err
		}
	}
	return nil
}

func (ws *Websocket) handleMessage(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty response")
	}

	switch data[0] {
	case '[':
		var events []WebsocketEvent
		if err := json.Unmarshal(data, &events); err != nil {
			return err
		}
		for i := range events {
			ws.events <- events[i]
		}
	case '{':
		var event WebsocketEvent
		if err := json.Unmarshal(data, &event); err == nil {
			ws.events <- event
			return nil
		}

		var resp WebsocketResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		if resp.Result != nil {
			return errors.Errorf("invalid websocket response: %v", resp.Result)
		}
	default:
		return errors.Errorf("invalid websocket response: %s", string(data))
	}

	return nil
}

func (ws *Websocket) send(req *WebsocketRequest) error {
	if ws.conn == nil {
		return nil
	}
	req.ID = ws.id
	ws.id++
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	ws.log.Trace().RawJSON("data", data).Msg("client->server")
	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

// SubscribeOnTickers -
func (ws *Websocket) SubscribeOnTickers(symbols ...string) error {
	params := make([]interface{}, 0)
	if len(symbols) == 0 {
		params = append(params, "!ticker@arr")
	} else {
		for i := range symbols {
			params = append(params, fmt.Sprintf("%s@ticker", strings.ToLower(symbols[i])))
		}
	}
	req := WebsocketRequest{
		Method: WebsocketMethodSubscribe,
		Params: params,
	}
	if err := ws.send(&req); err != nil {
		return err
	}
	ws.subscriptions[req.ID] = req
	return nil
}

// SubscribeOnBookTickers -
func (ws *Websocket) SubscribeOnBookTickers(symbols ...string) error {
	params := make([]interface{}, 0)
	if len(symbols) == 0 {
		params = append(params, "!bookTicker")
	} else {
		for i := range symbols {
			params = append(params, fmt.Sprintf("%s@bookTicker", strings.ToLower(symbols[i])))
		}
	}
	req := WebsocketRequest{
		Method: WebsocketMethodSubscribe,
		Params: params,
	}
	if err := ws.send(&req); err != nil {
		return err
	}
	ws.subscriptions[req.ID] = req
	return nil
}

// SubscribeOnKLine -
func (ws *Websocket) SubscribeOnKLine(interval Interval, symbol string) error {
	req := WebsocketRequest{
		Method: WebsocketMethodSubscribe,
		Params: []interface{}{
			fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval),
		},
	}
	if err := ws.send(&req); err != nil {
		return err
	}
	ws.subscriptions[req.ID] = req
	return nil
}
