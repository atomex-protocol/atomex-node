package main

import (
	"context"
	"flag"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var configDir string
	flag.StringVar(&configDir, "c", "configs", "path to directory containing configs")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	configDir = config.SelectEnvironment(configDir)

	var cfg Config
	configName := path.Join(configDir, "market_maker.yml")
	if err := config.Load(ctx, configName, &cfg); err != nil {
		log.Panic().Err(err).Str("file", configName).Msg("config.Load")
	}

	info, err := config.LoadGeneralConfig(ctx, configDir)
	if err != nil {
		log.Panic().Err(err).Msg("config.LoadGeneralConfig")
	}
	cfg.General = info

	if err := cfg.loadQuoteProviderMeta(ctx); err != nil {
		log.Panic().Err(err).Msg("config.loadQuoteProviderMeta")
	}

	marketMaker, err := NewMarketMaker(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("NewMarketMaker")
	}

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
	cancel()

	if err := marketMaker.Close(ctx); err != nil {
		log.Panic().Err(err).Msg("marketMaker.Close")
	}
	close(signals)

	log.Info().Msg("stopped")
}
