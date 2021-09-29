package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/rs/zerolog/log"
)

func run(cfg config.Config) error {
	watchTower, err := NewWatchTower(cfg)
	if err != nil {
		return err
	}
	if err := watchTower.Run(cfg.Restore); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals

	if err := watchTower.Close(); err != nil {
		return err
	}
	close(signals)

	log.Info().Msg("stopped")
	return nil
}
