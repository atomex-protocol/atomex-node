package config

import "github.com/atomex-protocol/watch_tower/internal/types"

// MetaInfo -
type MetaInfo struct {
	Assets  map[string]types.Asset
	Symbols []types.Symbol
}

// LoadMetaInfo -
func LoadMetaInfo(assetsFile, symbolsFile string) (info MetaInfo, err error) {
	info.Assets = make(map[string]types.Asset)
	if err = Load(assetsFile, &info.Assets); err != nil {
		return
	}

	info.Symbols = make([]types.Symbol, 0)
	if err = Load(symbolsFile, &info.Symbols); err != nil {
		return
	}

	for i := range info.Symbols {
		base := info.Symbols[i].BaseKey
		if asset, ok := info.Assets[base]; ok {
			info.Symbols[i].Base = asset
		}

		quote := info.Symbols[i].QuoteKey
		if asset, ok := info.Assets[quote]; ok {
			info.Symbols[i].Quote = asset
		}
	}

	return
}
