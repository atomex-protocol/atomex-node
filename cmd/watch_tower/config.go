package main

import (
	"github.com/atomex-protocol/watch_tower/internal/config"
)

// Config -
type Config struct {
	Restore              bool     `yaml:"restore"`
	Types                []string `yaml:"types" validate:"dive,oneof=redeem refund"`
	RetryCountOnFailedTx uint     `yaml:"retry_count_on_failed_tx"`

	General config.General `yaml:"-" validate:"-"`
}
