package atomex

import (
	"time"

	"github.com/rs/zerolog"
)

// RestOption -
type RestOption func(*Rest)

// WithTimeout -
func WithTimeout(timeout time.Duration) RestOption {
	return func(rest *Rest) {
		rest.timeout = timeout
	}
}

// WithURL -
func WithURL(baseURL string) RestOption {
	return func(rest *Rest) {
		rest.baseURL = baseURL
	}
}

// WithTestURL -
func WithTestURL() RestOption {
	return func(rest *Rest) {
		rest.baseURL = BaseTestURLRestAPI
	}
}

// WithProdURL -
func WithProdURL() RestOption {
	return func(rest *Rest) {
		rest.baseURL = BaseURLRestAPI
	}
}

// WithSignatureAlgorithm -
func WithSignatureAlgorithm(algo string) RestOption {
	return func(rest *Rest) {
		rest.algo = algo
	}
}

// WithRestLogLevel -
func WithRestLogLevel(level zerolog.Level) RestOption {
	return func(rest *Rest) {
		rest.log = rest.log.Level(level)
	}
}

// WebsocketOption -
type WebsocketOption func(*Websocket)

// WithWebsocketURI -
func WithWebsocketURI(uri string) WebsocketOption {
	return func(ws *Websocket) {
		ws.uri = uri
	}
}

// WithSignature -
func WithSignature(algo string) WebsocketOption {
	return func(ws *Websocket) {
		ws.algo = algo
	}
}

// WithLogLevel -
func WithLogLevel(level zerolog.Level) WebsocketOption {
	return func(ws *Websocket) {
		ws.logger = ws.logger.Level(level)
	}
}
