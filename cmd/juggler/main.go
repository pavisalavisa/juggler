package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pavisalavisa/juggler/internal/api"
	"github.com/pavisalavisa/juggler/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	configureLogger(cfg)

	run(cfg)
}

func run(cfg *config.Config) {
	ctx := context.Background()
	s := api.NewServer(cfg)

	log.Info().Msgf("Starting Juggler on port %d", cfg.Port)
	errors := s.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case err, ok := <-errors:
		if ok {
			log.Error().Err(err).Msgf("Juggler API server error")
		}
	case <-stop:
		log.Info().Msg("Shutdown signal received, shutting down")

		if err := s.Stop(ctx); err != nil {
			log.Fatal().Err(err).Msg("Error stopping server")
		}
	}

	log.Info().Msg("Shutdown successful!")
}

func configureLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		log.Warn().Err(err).Msg("cannot parse log level, using info")
		return
	}

	zerolog.SetGlobalLevel(level)
}
