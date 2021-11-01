package main

import (
	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/pkg/errors"
)

// Config -
type Config struct {
	Restore              bool     `yaml:"restore"`
	Types                []string `yaml:"types"`
	RetryCountOnFailedTx uint     `yaml:"retry_count_on_failed_tx"`

	Chains tools.Config `yaml:"-"`
}

// Validate -
func (c *Config) Validate() error {
	if err := c.Chains.Ethereum.Validate(); err != nil {
		return err
	}
	if err := c.Chains.Tezos.Validate(); err != nil {
		return err
	}
	for i := range c.Types {
		if c.Types[i] != "redeem" && c.Types[i] != "refund" {
			return errors.Errorf("invalid operation type (should be 'redeem' or 'refund')")
		}
	}
	return nil
}
