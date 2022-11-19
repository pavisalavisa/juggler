package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// Version is the application version (normally set when building the binary).
var Version = "dev"

type Config struct {
	// Port on which the juggler is running
	// Provided as JUGGLER_PORT
	Port int `default:"6060"`

	// Timeout to hard cancel of the proxy service after a SIGTERM
	//   provided as env variable JUGGLER_SHUTDOWN_TIMEOUT
	ShutdownTimeout time.Duration `default:"5s" split_words:"true"`

	// JUGGLER_READ_HEADER_TIMEOUT env var
	ReadHeaderTimeout time.Duration `default:"2s" envconfig:"READ_HEADER_TIMEOUT"`

	// Environment is the env where the service is running, provided by FIT.
	// env variable: FIT_ENV
	Environment string `default:"dev"`

	Logger Logger
}

type Logger struct {
	// Level defines the minimum level of severity that app should log.
	//   provided as env variable JUGGLER_LOGGER_LEVEL
	//   must be one of: ["trace", "debug", "info", "warn", "error", "critical"]
	Level string `default:"info" split_words:"true"`
}

func Load() *Config {
	prefix := "JUGGLER"

	cfg := &Config{}
	if err := envconfig.Process(prefix, cfg); err != nil {
		log.Fatal().Err(err).Msg("error loading configuration")
	}

	return cfg
}
