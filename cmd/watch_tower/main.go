package main

import (
	"flag"
	"os"
	"time"

	"github.com/aopoltorzhicky/watch_tower/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	var configName string
	flag.StringVar(&configName, "c", "config.yml", "path to YAML config file")

	flag.Parse()

	cfg, err := config.Load(configName)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := cfg.Validate(); err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := run(cfg); err != nil {
		log.Panic().Stack().Err(err).Msg("")
	}
}
