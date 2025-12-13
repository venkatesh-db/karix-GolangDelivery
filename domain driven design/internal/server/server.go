package server

import (
	"context"
	"net/http"
	"time"

	"github.com/helrachar/banking/internal/config"
)

// HTTPServer wraps net/http server configuration and lifecycle helpers.
type HTTPServer struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

// NewHTTPServer builds an HTTP server with the provided configuration and handler.
func NewHTTPServer(cfg config.HTTPConfig, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

// Start begins serving HTTP requests.
func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown attempts a graceful shutdown within the provided context.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// ShutdownTimeout exposes the configured graceful shutdown window.
func (s *HTTPServer) ShutdownTimeout() time.Duration {
	return s.shutdownTimeout
}
