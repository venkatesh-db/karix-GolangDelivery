package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"go.uber.org/zap"

	"project_saas/shared/pkg/config"
)

// Run spins up an HTTP server with graceful shutdown.
func Run(ctx context.Context, cfg config.ServiceConfig, log *zap.Logger, register func(r chi.Router)) error {
	if register == nil {
		return fmt.Errorf("register func required")
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(RequestLogger(log))
	r.Use(otelchi.Middleware(cfg.ServiceName, otelchi.WithChiRoutes(r)))
	r.Handle("/metrics", promhttp.Handler())
	register(r)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("http server starting", zap.String("port", cfg.HTTPPort))
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("http server shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
