package main

import (
	"context"
	"fmt"
	"path"

	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/synthetic"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/atomex-protocol/watch_tower/internal/keys"
)

// Config -
type Config struct {
	QuoteProvider QuoteProvider     `yaml:"quote_provider" validate:"required"`
	Strategies    []strategy.Config `yaml:"strategies" validate:"required"`
	Keys          Keys              `yaml:"keys" validate:"required"`
	LogLevel      string            `yaml:"log_level"`
	Restore       bool              `yaml:"restore"`

	General           config.General    `yaml:"-" validate:"-"`
	QuoteProviderMeta QuoteProviderMeta `yaml:"-" validate:"-"`
}

func (cfg *Config) loadQuoteProviderMeta(ctx context.Context) error {
	quoteProviderFile := path.Join(cfg.General.ConfigDir, fmt.Sprintf("%s.yml", cfg.QuoteProvider.Kind))
	var quoteProviderConfig QuoteProviderMeta
	if err := config.Load(ctx, quoteProviderFile, &quoteProviderConfig); err != nil {
		return err
	}

	cfg.QuoteProviderMeta = quoteProviderConfig
	return nil
}

// QuoteProvider -
type QuoteProvider struct {
	Kind QuoteProviderKind `yaml:"kind" validate:"required,oneof=binance"`
}

// QuoteProviderKind -
type QuoteProviderKind string

// quote provider kinds
const (
	QuoteProviderKindBinance QuoteProviderKind = "binance"
)

// QuoteProviderMeta -
type QuoteProviderMeta struct {
	FromSymbols map[string]synthetic.Config `yaml:"from_symbols" validate:"required"`
	ToSymbols   map[string]string           `yaml:"to_symbols" validate:"required"`
}

// Keys -
type Keys struct {
	Kind                keys.StorageKind `yaml:"kind"`
	File                string           `yaml:"file"`
	GenerateIfNotExists bool             `yaml:"generate_if_not_exists"`
}
