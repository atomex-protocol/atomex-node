package tools

import "errors"

const (
	chainsCount = 2
)

// Config -
type Config struct {
	Tezos    Tezos    `yaml:"tezos"`
	Ethereum Ethereum `yaml:"ethereum"`
}

// Tezos -
type Tezos struct {
	MinPayOff string   `yaml:"min_payoff"`
	Node      string   `yaml:"node"`
	TzKT      string   `yaml:"tzkt"`
	Tokens    []string `yaml:"tokens"`
	Contract  string   `yaml:"contract"`
	TTL       int64    `yaml:"ttl"`
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
	if t.TTL == 0 {
		t.TTL = 5
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
	if e.MinPayOff == "" {
		e.MinPayOff = "0"
	}
	return nil
}
