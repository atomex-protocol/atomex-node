package main

import (
	"context"
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

	ctx := context.Background()

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

	if err := run(cfg); err != nil {
		log.Panic().Stack().Err(err).Msg("")
	}
}
