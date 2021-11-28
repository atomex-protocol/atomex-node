package atomex

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/atomex/signers"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Rest - realization of Atomex REST API
type Rest struct {
	baseURL string
	token   string
	algo    string
	timeout time.Duration

	log zerolog.Logger
}

// NewRest - constructor of `Rest`
func NewRest(opts ...RestOption) *Rest {
	r := &Rest{
		log: logger.New(logger.WithModuleName("atomex_rest_api")),
	}
	for i := range opts {
		opts[i](r)
	}
	if r.baseURL == "" {
		r.baseURL = BaseURLRestAPIv1
	}
	if r.timeout < 1 {
		r.timeout = time.Second * 10
	}
	if r.algo == "" {
		r.algo = signers.AlgorithmBlake2bWithEcdsaSecp256k1
	}
	return r
}

func (rest *Rest) request(ctx context.Context, method string, path string, args url.Values, body interface{}, output interface{}) error {
	client := http.Client{
		Timeout: rest.timeout,
	}

	uri, err := url.Parse(fmt.Sprintf("%s/%s", rest.baseURL, path))
	if err != nil {
		return err
	}
	if len(args) > 0 {
		uri.RawQuery = args.Encode()
	}

	trace := rest.log.Trace().
		Str("uri", uri.String()).
		Str("method", method)

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "body")
		}
		bodyReader = bytes.NewBuffer(data)

		trace = trace.RawJSON("body", data)
	}

	trace.Msg("client->server")

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), bodyReader)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	if rest.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rest.token))
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if output != nil {
			return json.NewDecoder(resp.Body).Decode(&output)
		}
		return nil
	default:
		var atomexErr Error
		if err := json.NewDecoder(resp.Body).Decode(&atomexErr); err != nil {
			return errors.Wrap(err, "error decoding in atomex request")
		}
		return errors.Wrapf(atomexErr, "request (status code: %d)", resp.StatusCode)
	}
}

// Token - get authentication token
func (rest *Rest) Token(ctx context.Context, keys *signers.Key) (response TokenResponse, err error) {
	req := NewTokenRequest(signMessage, rest.algo, keys.Public)
	if err := req.Sign(keys); err != nil {
		return response, errors.Wrap(err, "sign")
	}

	err = rest.request(ctx, http.MethodPost, "Token", nil, req, &response)
	return
}

// Auth - authenticate `Rest` in atomex server
func (rest *Rest) Auth(ctx context.Context, keys *signers.Key) error {
	token, err := rest.Token(ctx, keys)
	if err != nil {
		return errors.Wrap(err, "Token")
	}
	rest.token = token.Token
	return nil
}

// GetToken -
func (rest *Rest) GetToken() string {
	return rest.token
}

// TopOfBookQuotes -
func (rest *Rest) TopOfBookQuotes(ctx context.Context, symbols ...string) (response []TopOfBook, err error) {
	args := make(url.Values)
	if len(symbols) > 0 {
		args.Add("symbols", strings.Join(symbols, ","))
	}
	err = rest.request(ctx, http.MethodGet, "MarketData/quotes", args, nil, &response)
	return
}

// OrderBook -
func (rest *Rest) OrderBook(ctx context.Context, symbol string) (response OrderBook, err error) {
	if symbol == "" {
		return response, errors.New("empty symbol in OrderBook request")
	}
	args := make(url.Values)
	args.Add("symbol", symbol)
	err = rest.request(ctx, http.MethodGet, "MarketData/book", args, nil, &response)
	return
}

// AddOrder -
func (rest *Rest) AddOrder(ctx context.Context, req AddOrderRequest) (response AddOrderResponse, err error) {
	err = rest.request(ctx, http.MethodPost, "Orders", nil, req, &response)
	return
}

// Orders -
func (rest *Rest) Orders(ctx context.Context, req OrdersRequest) (response []Order, err error) {
	err = rest.request(ctx, http.MethodGet, "Orders", req.getArgs(), nil, &response)
	return
}

// Order -
func (rest *Rest) Order(ctx context.Context, id int64) (response Order, err error) {
	if id < 1 {
		return response, errors.Errorf("invalid order id: %d", id)
	}
	urlPath := fmt.Sprintf("Orders/%d", id)
	err = rest.request(ctx, http.MethodGet, urlPath, nil, nil, &response)
	return
}

// CancelOrder -
func (rest *Rest) CancelOrder(ctx context.Context, id int64, symbol string, side Side) (response DefaultResponse, err error) {
	if id < 1 {
		return response, errors.Errorf("invalid order id: %d", id)
	}

	args := make(url.Values)
	if symbol != "" {
		args.Add("symbol", symbol)
	}
	if side != SideEmpty {
		args.Add("side", string(side))
	}
	urlPath := fmt.Sprintf("Orders/%d", id)
	err = rest.request(ctx, http.MethodDelete, urlPath, args, nil, &response)
	return
}

// Swap -
func (rest *Rest) Swap(ctx context.Context, id int64) (response Swap, err error) {
	if id < 1 {
		return response, errors.Errorf("invalid order id: %d", id)
	}
	urlPath := fmt.Sprintf("Swaps/%d", id)
	err = rest.request(ctx, http.MethodGet, urlPath, nil, nil, &response)
	return
}

// Swaps -
func (rest *Rest) Swaps(ctx context.Context, req SwapsRequest) (response []Swap, err error) {
	err = rest.request(ctx, http.MethodGet, "Swaps", req.getArgs(), nil, &response)
	return
}

// AddSwapRequisites -
func (rest *Rest) AddSwapRequisites(ctx context.Context, id int64, req AddSwapRequisitesRequest) (response DefaultResponse, err error) {
	urlPath := fmt.Sprintf("Swaps/%d/requisites", id)
	err = rest.request(ctx, http.MethodPost, urlPath, nil, req, &response)
	return
}

// SymbolInfo -
func (rest *Rest) SymbolInfo(ctx context.Context) (response []SymbolInfo, err error) {
	err = rest.request(ctx, http.MethodGet, "Symbols", nil, nil, &response)
	return
}
