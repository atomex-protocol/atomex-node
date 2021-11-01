package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New -
func New(opts ...LoggerOption) zerolog.Logger {
	args := newArgs()
	for i := range opts {
		opts[i](&args)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := log.
		Output(output).
		Level(args.Level).
		With().
		Timestamp()

	if args.Module != "" {
		logger = logger.Str("module", args.Module)
	}

	return logger.Logger()
}

type args struct {
	Level  zerolog.Level
	Module string
}

func newArgs() args {
	return args{
		Level: zerolog.InfoLevel,
	}
}

// LoggerOption -
type LoggerOption func(*args)

// WithLogLevel -
func WithLogLevel(level zerolog.Level) LoggerOption {
	return func(ws *args) {
		ws.Level = level
	}
}

// WithModuleName -
func WithModuleName(module string) LoggerOption {
	return func(ws *args) {
		ws.Module = module
	}
}
