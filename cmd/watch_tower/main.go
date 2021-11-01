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

	var configName string
	flag.StringVar(&configName, "c", "config.yml", "path to YAML config file")

	flag.Parse()

	var cfg Config
	if err := config.Load(configName, &cfg); err != nil {
		log.Panic().Err(err).Msg("")
	}

	var configDir string
	env := os.Getenv("ATOMEX_PROTOCOL_ENV")
	switch env {
	case "production":
		configDir = "./configs"
	case "test":
		configDir = "../../configs/test"
	default:
		log.Panic().Str("env", env).Msg("invalid environment")
	}

	if err := config.Load(path.Join(configDir, "chains.yml"), &cfg.Chains); err != nil {
		log.Panic().Err(err).Str("file", path.Join(configDir, "chains.yml")).Msg("config.Load")
	}

	if err := cfg.Validate(); err != nil {
		log.Panic().Err(err).Msg("")
	}

	if err := run(cfg); err != nil {
		log.Panic().Stack().Err(err).Msg("")
	}
}
