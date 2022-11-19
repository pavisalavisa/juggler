package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/pavisalavisa/juggler/internal/api"
	"github.com/pavisalavisa/juggler/internal/config"

	"github.com/stretchr/testify/require"
)

func TestServer_NewServer_Running(t *testing.T) {
	cfg := fixtureConfig()

	server := api.NewServer(cfg)

	errChan := server.Start()
	err := server.Stop(context.Background())

	require.NoError(t, err, "Stopping a server should not return an error")
	err = <-errChan
	require.ErrorIs(t, err, http.ErrServerClosed, "Server closed should be returned by server channel")
}

func fixtureConfig() *config.Config {
	return &config.Config{
		Port:              6060,
		ShutdownTimeout:   10,
		Environment:       "test",
		Logger:            config.Logger{Level: "info"},
		ReadHeaderTimeout: 500}
}
