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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var configDir string
	flag.StringVar(&configDir, "c", "configs", "path to configs directory")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	configDir = config.SelectEnvironment(configDir)

	var cfg Config
	configName := path.Join(configDir, "watch_tower.yml")
	if err := config.Load(ctx, configName, &cfg); err != nil {
		log.Panic().Err(err).Msg("config.Load")
	}

	general, err := config.LoadGeneralConfig(ctx, configDir)
	if err != nil {
		log.Panic().Err(err).Msg("LoadGeneralConfig")
	}
	cfg.General = general

	watchTower, err := NewWatchTower(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("NewWatchTower")
	}

	if err := watchTower.Run(ctx, cfg.Restore); err != nil {
		log.Panic().Err(err).Msg("Run")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	cancel()

	if err := watchTower.Close(); err != nil {
		log.Panic().Err(err).Msg("Close")
	}
	close(signals)

	log.Info().Msg("stopped")
}
