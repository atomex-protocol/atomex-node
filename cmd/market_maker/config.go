package main

import (
	"github.com/atomex-protocol/watch_tower/cmd/market_maker/strategy"
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/atomex-protocol/watch_tower/internal/keys"
)

// Config -
type Config struct {
	QuoteProvider QuoteProvider     `yaml:"quote_provider"`
	Strategies    []strategy.Config `yaml:"strategies"`
	Keys          Keys              `yaml:"keys"`
	LogLevel      string            `yaml:"log_level"`

	Info   config.MetaInfo `yaml:"-"`
	Atomex config.Atomex   `yaml:"-"`
	Chains tools.Config    `yaml:"-"`
}

// QuoteProvider -
type QuoteProvider struct {
	Kind QuoteProviderKind `yaml:"kind"`

	Meta QuoteProviderMeta `yaml:"-"`
}

// QuoteProviderKind -
type QuoteProviderKind string

// quote provider kinds
const (
	QuoteProviderKindBinance QuoteProviderKind = "binance"
)

// QuoteProviderMeta -
type QuoteProviderMeta struct {
	FromSymbols map[string]string `yaml:"from_symbols"`
	ToSymbols   map[string]string `yaml:"to_symbols"`
}

// Keys -
type Keys struct {
	Kind                keys.StorageKind `yaml:"kind"`
	File                string           `yaml:"file"`
	GenerateIfNotExists bool             `yaml:"generate_if_not_exists"`
}
