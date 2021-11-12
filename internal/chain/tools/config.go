package tools

import (
	"errors"

	"github.com/atomex-protocol/watch_tower/internal/chain/tezos"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/go-playground/validator/v10"
)

const (
	chainsCount = 2
)

// Config -
type Config struct {
	Tezos    Tezos    `yaml:"tezos" validate:"required"`
	Ethereum Ethereum `yaml:"ethereum" validate:"required"`
}

// Tezos -
type Tezos struct {
	MinPayOff string `yaml:"min_payoff"`
	Node      string `yaml:"node" validate:"required,uri"`
	TzKT      string `yaml:"tzkt" validate:"required,uri"`
	TTL       int64  `yaml:"ttl"`

	Tokens          []string                         `yaml:"-" validate:"-"`
	Contract        string                           `yaml:"-" validate:"-"`
	OperaitonParams tezos.OperationParamsByContracts `yaml:"-" validate:"-"`
}

// FillContractAddresses -
func (t *Tezos) FillContractAddresses(assets map[string]types.Asset) error {
	xtz, ok := assets["XTZ"]
	if !ok {
		return errors.New("assets.yml does not contains XTZ asset")
	}
	t.Contract = xtz.AtomexContract

	tokens := make(map[string]struct{})
	for name, asset := range assets {
		if asset.Chain != "tezos" || name == "XTZ" || asset.AtomexContract == "" {
			continue
		}

		tokens[asset.AtomexContract] = struct{}{}
	}

	t.Tokens = make([]string, 0)
	for address := range tokens {
		t.Tokens = append(t.Tokens, address)
	}
	return nil
}

// Validate -
func (t *Tezos) Validate() error {
	if err := validator.New().Struct(t); err != nil {
		return err
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
	Node         string `yaml:"node" validate:"required,uri"`
	Wss          string `yaml:"wss" validate:"required,uri"`
	EthAddress   string `yaml:"-" validate:"-"`
	Erc20Address string `yaml:"-" validate:"-"`
}

// Validate -
func (e *Ethereum) Validate() error {
	if err := validator.New().Struct(e); err != nil {
		return err
	}
	if e.MinPayOff == "" {
		e.MinPayOff = "0"
	}
	return nil
}

// FillContractAddresses -
func (e *Ethereum) FillContractAddresses(assets map[string]types.Asset) error {
	xtz, ok := assets["ETH"]
	if !ok {
		return errors.New("assets.yml does not contains ETH asset")
	}
	e.EthAddress = xtz.AtomexContract

	for name, asset := range assets {
		if asset.Chain != "ethereum" || name == "ETH" {
			continue
		}

		if asset.AtomexContract != "" {
			e.Erc20Address = asset.AtomexContract
		}
	}

	return nil
}
