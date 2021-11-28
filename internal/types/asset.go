package types

import "github.com/atomex-protocol/watch_tower/internal/chain"

// Asset -
type Asset struct {
	Name           string `yaml:"name" validate:"require"`
	Chain          string `yaml:"chain" validate:"require"`
	Contract       string `yaml:"contract"`
	AtomexContract string `yaml:"atomex_contract" validate:"require"`
	Decimals       int    `yaml:"decimals" validate:"require"`
}

// ChainType -
func (a Asset) ChainType() chain.ChainType {
	switch a.Chain {
	case "ethereum":
		return chain.ChainTypeEthereum
	case "tezos":
		return chain.ChainTypeTezos
	default:
		return chain.ChainTypeUnknown
	}
}
