package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var configDir string
	flag.StringVar(&configDir, "c", "configs", "path to directory containing configs")

	flag.Parse()

	configDir = config.SelectEnvironment(configDir)

	var cfg Config
	configName := path.Join(configDir, "market_maker.yml")
	if err := config.Load(configName, &cfg); err != nil {
		log.Panic().Err(err).Str("file", configName).Msg("config.Load")
	}

	info, err := config.LoadMetaInfo(path.Join(configDir, "assets.yml"), path.Join(configDir, "symbols.yml"))
	if err != nil {
		log.Panic().Err(err).Msg("config.LoadMetaInfo")
	}
	cfg.Info = info

	if err := config.Load(path.Join(configDir, "atomex.yml"), &cfg.Atomex); err != nil {
		log.Panic().Err(err).Str("file", path.Join(configDir, "atomex.yml")).Msg("config.Load")
	}

	if err := config.Load(path.Join(configDir, "chains.yml"), &cfg.Chains); err != nil {
		log.Panic().Err(err).Str("file", path.Join(configDir, "chains.yml")).Msg("config.Load")
	}

	if err := cfg.Chains.Ethereum.FillContractAddresses(cfg.Info.Assets); err != nil {
		log.Panic().Err(err).Msg("FillContractAddresses")
	}
	if err := cfg.Chains.Tezos.FillContractAddresses(cfg.Info.Assets); err != nil {
		log.Panic().Err(err).Msg("FillContractAddresses")
	}

	if err := cfg.Chains.Ethereum.Validate(); err != nil {
		log.Panic().Err(err).Msg("config.Ethereum.Validate")
	}

	if err := cfg.Chains.Tezos.Validate(); err != nil {
		log.Panic().Err(err).Msg("config.Tezos.Validate")
	}

	quoteProviderFile := path.Join(configDir, fmt.Sprintf("%s.yml", cfg.QuoteProvider.Kind))
	var quoteProviderConfig QuoteProviderMeta
	if err := config.Load(quoteProviderFile, &quoteProviderConfig); err != nil {
		log.Panic().Err(err).Str("file", quoteProviderFile).Msg("config.Load")
	}
	cfg.QuoteProvider.Meta = quoteProviderConfig

	marketMaker, err := NewMarketMaker(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("NewMarketMaker")
	}

	ctx := context.Background()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("panic", err).Msg("panic occurred")

			signals <- syscall.SIGINT
		}
	}()

	if err := marketMaker.Start(ctx); err != nil {
		log.Panic().Err(err).Msg("marketMaker.Start")
	}

	<-signals

	if err := marketMaker.Close(ctx); err != nil {
		log.Panic().Err(err).Msg("marketMaker.Close")
	}
	close(signals)

	log.Info().Msg("stopped")
}
