package main

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var configDir string
	flag.StringVar(&configDir, "c", "configs", "path to configs directory")

	flag.Parse()

	configDir = config.SelectEnvironment(configDir)

	var cfg Config
	configName := path.Join(configDir, "watch_tower.yml")
	if err := config.Load(configName, &cfg); err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := config.Load(path.Join(configDir, "chains.yml"), &cfg.Chains); err != nil {
		log.Panic().Err(err).Str("file", path.Join(configDir, "chains.yml")).Msg("config.Load")
	}

	if err := config.Load(path.Join(configDir, "assets.yml"), &cfg.Assets); err != nil {
		log.Panic().Err(err).Str("file", path.Join(configDir, "chains.yml")).Msg("config.Load")
	}

	if err := cfg.Chains.Ethereum.FillContractAddresses(cfg.Assets); err != nil {
		log.Panic().Err(err).Msg("FillContractAddresses")
	}
	if err := cfg.Chains.Tezos.FillContractAddresses(cfg.Assets); err != nil {
		log.Panic().Err(err).Msg("FillContractAddresses")
	}

	if err := cfg.Validate(); err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := run(cfg); err != nil {
		log.Panic().Stack().Err(err).Msg("")
	}
}
