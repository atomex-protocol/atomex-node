package config

import (
	"context"
	"path"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain/tools"
	"github.com/atomex-protocol/watch_tower/internal/types"
	"github.com/pkg/errors"
)

// General -
type General struct {
	Assets  Assets
	Symbols []types.Symbol
	Chains  tools.Config
	Atomex  Atomex

	ConfigDir string
}

// Assets
type Assets map[string]types.Asset

// LoadGeneralConfig -
func LoadGeneralConfig(ctx context.Context, configDir string) (general General, err error) {
	general.ConfigDir = configDir

	symbols, err := loadSymbols(ctx, configDir)
	if err != nil {
		return
	}
	general.Symbols = symbols

	assets, err := loadAssets(ctx, configDir)
	if err != nil {
		return
	}
	general.Assets = assets

	for i := range general.Symbols {
		base := general.Symbols[i].BaseKey
		if asset, ok := general.Assets[base]; ok {
			general.Symbols[i].Base = asset
		}

		quote := general.Symbols[i].QuoteKey
		if asset, ok := general.Assets[quote]; ok {
			general.Symbols[i].Quote = asset
		}
	}

	chains, err := loadChains(ctx, configDir)
	if err != nil {
		return
	}
	general.Chains = chains

	if err := general.Chains.Ethereum.FillContractAddresses(general.Assets); err != nil {
		return general, errors.Wrap(err, "Ethereum.FillContractAddresses")
	}
	if err := general.Chains.Tezos.FillContractAddresses(general.Assets); err != nil {
		return general, errors.Wrap(err, "Tezos.FillContractAddresses")
	}

	atomexConfig, err := loadAtomex(ctx, configDir)
	if err != nil {
		return
	}
	general.Atomex = atomexConfig

	return
}

func loadSymbols(ctx context.Context, configDir string) ([]types.Symbol, error) {
	internalCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	symbolsFile := path.Join(configDir, "symbols.yml")
	symbols := make([]types.Symbol, 0)
	err := Load(internalCtx, symbolsFile, &symbols)
	return symbols, err
}

func loadAssets(ctx context.Context, configDir string) (Assets, error) {
	// TODO: assets verification
	internalCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	assetsFile := path.Join(configDir, "assets.yml")
	assets := make(Assets)
	err := Load(internalCtx, assetsFile, &assets)
	return assets, err
}

func loadChains(ctx context.Context, configDir string) (tools.Config, error) {
	internalCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var chains tools.Config
	chainsFile := path.Join(configDir, "chains.yml")
	if err := Load(internalCtx, chainsFile, &chains); err != nil {
		return chains, err
	}

	tezosCtx, cancelTezos := context.WithTimeout(ctx, time.Second)
	defer cancelTezos()
	err := Load(tezosCtx, path.Join(configDir, "tezos.yml"), &chains.Tezos.OperaitonParams)
	return chains, err
}

func loadAtomex(ctx context.Context, configDir string) (Atomex, error) {
	internalCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var cfg Atomex
	atomexFile := path.Join(configDir, "atomex.yml")
	err := Load(internalCtx, atomexFile, &cfg)
	return cfg, err
}
