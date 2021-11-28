package binance

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/exchange"
	"github.com/atomex-protocol/watch_tower/internal/secrets"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// Rest -
type Rest struct {
	url           string
	delta         int64
	publicLimiter *rate.Limiter
	log           zerolog.Logger
}

func newRest(url string, logger zerolog.Logger) *Rest {
	return &Rest{
		url:           url,
		publicLimiter: rate.NewLimiter(rate.Every(time.Minute), 1200),
		log:           logger,
	}
}

// Init -
func (rest *Rest) Init() error {
	serverTime, err := rest.ServerTime()
	if err != nil {
		return err
	}
	rest.delta = time.Now().UnixNano()/1_000_000 - serverTime

	info, err := rest.ExhangeInfo()
	if err != nil {
		return err
	}

	for _, limit := range info.RateLimits {
		if limit.RateLimitType == "REQUEST_WEIGHT" {
			rest.publicLimiter = getLimitInterval(limit)
		}
	}
	return nil
}

func getLimitInterval(limit RateLimit) *rate.Limiter {
	var duration time.Duration
	switch limit.Interval {
	case "SECOND":
		duration = time.Second
	case "MINUTE":
		duration = time.Minute
	case "HOUR":
		duration = time.Hour
	case "DAY":
		duration = time.Hour * 24
	default:
		return nil
	}
	return rate.NewLimiter(rate.Every(duration*time.Duration(limit.IntervalNum)), int(limit.Limit))
}

func (rest *Rest) request(isPrivate bool, method, path string, args url.Values, body url.Values, weight int, output interface{}) error {
	uri, err := url.Parse(fmt.Sprintf("%s/%s", rest.url, path))
	if err != nil {
		return err
	}
	if isPrivate {
		args.Add("timestamp", fmt.Sprintf("%v", time.Now().UnixNano()/1_000_000-rest.delta))
	}

	if len(args) > 0 {
		uri.RawQuery = args.Encode()
	}

	if err := rest.publicLimiter.WaitN(context.Background(), weight); err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri.String(), bytes.NewBufferString(body.Encode()))
	if err != nil {
		return err
	}

	if isPrivate {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err := rest.auth(req, args, body); err != nil {
			return err
		}
	}

	rest.log.Trace().
		Str("url", uri.String()).
		Str("method", method).
		Str("body", body.Encode()).
		Msg("request")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		if output != nil {
			return json.NewDecoder(response.Body).Decode(output)
		}
		return nil
	case http.StatusTooManyRequests, http.StatusTeapot:
		retryAfter := time.Now()
		if retryAfterStr := response.Header.Get("Retry-After"); retryAfterStr != "" {
			if seconds, err := strconv.Atoi(retryAfterStr); err == nil {
				retryAfter = retryAfter.Add(time.Second * time.Duration(seconds))
			} else {
				retryAfter = retryAfter.Add(time.Second * 10)
			}
		} else {
			retryAfter = retryAfter.Add(time.Second * 10)
		}
		return exchange.ErrToManyRequests{
			RetryAfter: retryAfter,
		}
	default:
		resp, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var binanceErr Error
		if err := json.Unmarshal(resp, &binanceErr); err != nil {
			return errors.Errorf("invalid status: %s (%s)", response.Status, string(resp))
		}
		return errors.Wrap(binanceErr.Handle(), uri.String())
	}
}

func (rest *Rest) auth(req *http.Request, args, body url.Values) error {
	req.Header.Add("X-MBX-APIKEY", secrets.Load("BINANCE_API_KEY"))

	raw := fmt.Sprintf("%s%s", args.Encode(), body.Encode())
	mac := hmac.New(sha256.New, []byte(secrets.Load("BINANCE_API_SECRET")))
	if _, err := mac.Write([]byte(raw)); err != nil {
		return err
	}
	v := url.Values{}
	v.Set("signature", fmt.Sprintf("%x", (mac.Sum(nil))))
	if len(args) == 0 {
		req.URL.RawQuery = v.Encode()
	} else {
		req.URL.RawQuery += fmt.Sprintf("&%s", v.Encode())
	}
	return nil
}

// ServerTime - returns server time in milliseconds
func (rest *Rest) ServerTime() (int64, error) {
	var data Server
	if err := rest.request(false, http.MethodGet, "api/v3/time", url.Values{}, url.Values{}, 1, &data); err != nil {
		return 0, err
	}
	return data.Time, nil
}

// ExhangeInfo -
func (rest *Rest) ExhangeInfo(symbols ...string) (data Info, err error) {
	args := url.Values{}
	switch len(symbols) {
	case 0:
	case 1:
		args.Add("symbol", symbols[0])
	default:
		args.Add("symbols", fmt.Sprintf("%v", symbols))
	}
	err = rest.request(false, http.MethodGet, "api/v3/exchangeInfo", args, url.Values{}, 10, &data)
	return
}

// ExhangeInfo -
func (rest *Rest) OHLC(symbol string, interval Interval, start, end, limit uint64) (data []OHLC, err error) {
	args := url.Values{}
	args.Add("symbol", symbol)
	args.Add("interval", string(interval))
	if start > 0 {
		args.Add("startTime", strconv.FormatUint(start, 10))
	}
	if end > 0 {
		args.Add("endTime", strconv.FormatUint(end, 10))
	}
	if limit > 0 && limit <= 1000 {
		args.Add("limit", strconv.FormatUint(limit, 10))
	}

	err = rest.request(false, http.MethodGet, "api/v3/klines", args, url.Values{}, 1, &data)
	return
}
