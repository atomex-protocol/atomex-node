package tools

import "github.com/rs/zerolog"

// TrackerOption -
type TrackerOption func(*Tracker)

// WithLogLevel -
func WithLogLevel(level zerolog.Level) TrackerOption {
	return func(t *Tracker) {
		t.logger = t.logger.Level(level)
	}
}

// WithRestore -
func WithRestore() TrackerOption {
	return func(t *Tracker) {
		t.needRestore = true
	}
}
