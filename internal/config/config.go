package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	Tezos                Tezos    `yaml:"tezos"`
	Ethereum             Ethereum `yaml:"ethereum"`
	Restore              bool     `yaml:"restore"`
	Types                []string `yaml:"types"`
	RetryCountOnFailedTx uint     `yaml:"retry_count_on_failed_tx"`
}

// Validate -
func (c *Config) Validate() error {
	if err := c.Ethereum.Validate(); err != nil {
		return err
	}
	if err := c.Tezos.Validate(); err != nil {
		return err
	}
	for i := range c.Types {
		if c.Types[i] != "redeem" && c.Types[i] != "refund" {
			return errors.Errorf("invalid operation type (should be 'redeem' or 'refund')")
		}
	}
	return nil
}

// Tezos -
type Tezos struct {
	MinPayOff string   `yaml:"min_payoff"`
	Node      string   `yaml:"node"`
	TzKT      string   `yaml:"tzkt"`
	Tokens    []string `yaml:"tokens"`
	Contract  string   `yaml:"contract"`
}

// Validate -
func (t *Tezos) Validate() error {
	if t.Node == "" {
		return errors.New("empty tezos node URL (key `tezos.node`)")
	}
	if t.TzKT == "" {
		return errors.New("empty tzkt URL (key `tezos.tzkt`)")
	}
	if t.Contract == "" {
		return errors.New("empty atomex tezos contract (key `tezos.contract`)")
	}
	if len(t.Tokens) == 0 {
		return errors.New("empty atomex tezos token contracts (key `tezos.tokens`)")
	}
	if t.MinPayOff == "" {
		t.MinPayOff = "0"
	}
	return nil
}

// Ethereum -
type Ethereum struct {
	MinPayOff    string `yaml:"min_payoff"`
	Node         string `yaml:"node"`
	Wss          string `yaml:"wss"`
	EthAddress   string `yaml:"eth_address"`
	Erc20Address string `yaml:"erc20_address"`
	UserAddress  string `yaml:"user_address"`
}

// Validate -
func (e *Ethereum) Validate() error {
	if e.Node == "" {
		return errors.New("empty ethereum node URL (key `ethereum.node`)")
	}
	if e.Wss == "" {
		return errors.New("empty ethereum websocket URL (key `ethereum.wss`)")
	}
	if e.EthAddress == "" {
		return errors.New("empty atomex ethereum contract address (key `ethereum.eth_address`)")
	}
	if e.Erc20Address == "" {
		return errors.New("empty atomex erc20 contract address (key `ethereum.erc20_address`)")
	}
	if e.UserAddress == "" {
		return errors.New("empty ethereum user wallet address address (key `ethereum.user_address`)")
	}
	if e.MinPayOff == "" {
		e.MinPayOff = "0"
	}
	return nil
}

// Load -
func Load(filename string) (c Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&c)
	return c, err
}
