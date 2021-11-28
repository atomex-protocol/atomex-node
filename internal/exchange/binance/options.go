package binance

import "github.com/rs/zerolog"

type options struct {
	Level       zerolog.Level
	BaseURLRest string
	BaseURLWs   string
}

func newOptions() options {
	return options{
		Level:       zerolog.InfoLevel,
		BaseURLRest: BaseURLServer2,
		BaseURLWs:   BaseURLWebsocket,
	}
}

// BinanceOption -
type BinanceOption func(*options)

// WithRestURL -
func WithRestURL(url string) BinanceOption {
	return func(opt *options) {
		opt.BaseURLRest = url
	}
}

// WithWebsocketURL -
func WithWebsocketURL(url string) BinanceOption {
	return func(opt *options) {
		opt.BaseURLWs = url
	}
}

// WithLogLevel -
func WithLogLevel(level zerolog.Level) BinanceOption {
	return func(opt *options) {
		opt.Level = level
	}
}
