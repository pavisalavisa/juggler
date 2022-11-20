package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pavisalavisa/juggler/internal/config"
	"github.com/pavisalavisa/juggler/internal/proxy"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type JugglerServer struct {
	cfg    config.Config
	server http.Server
}

func NewServer(cfg *config.Config) *JugglerServer {
	s := &JugglerServer{cfg: *cfg}

	safetyTimeout := 100 * time.Millisecond // give tcp more slack and use a Timeout middleware instead

	s.server = http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           s.router(cfg),
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout + safetyTimeout,
		ReadTimeout:       2*cfg.ReadHeaderTimeout + safetyTimeout,
		WriteTimeout:      2*cfg.ReadHeaderTimeout + safetyTimeout,
	}

	return s
}

func (s *JugglerServer) Start() <-chan error {
	proxyerrors := make(chan error)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			proxyerrors <- fmt.Errorf("error starting server: %w", err)
		}

		close(proxyerrors)
	}()

	return proxyerrors
}

func (s *JugglerServer) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	log.Info().Msg("Shutting down juggler server")
	if err := s.server.Shutdown(ctx); err != nil {
		_ = s.server.Close()
		return err
	}

	return nil
}

func (s *JugglerServer) router(cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	withLogger(r)
	r.Get("/status", s.statusHandler)

	proxy := NewProxyService(proxy.NewProxy(proxy.NoOpCaller{}))

	r.Handle("/", proxy.proxyHandler())
	return r
}

func (s *JugglerServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK!\n"))
}

func withLogger(r *chi.Mux) {
	r.Use(
		hlog.NewHandler(log.Logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
		hlog.RequestIDHandler("req_id", "Request-Id"))
}
