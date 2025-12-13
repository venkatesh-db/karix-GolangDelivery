package bootstrap

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"project_saas/shared/pkg/config"
	"project_saas/shared/pkg/httpx"
	"project_saas/shared/pkg/logger"
	"project_saas/shared/pkg/observability"
)

// RunHTTPService wires up config, logging, and HTTP server launch for a service.
func RunHTTPService(ctx context.Context, serviceName string, register func(r chi.Router, cfg config.ServiceConfig, log *zap.Logger)) error {
	cfg, err := config.Load(serviceName)
	if err != nil {
		return err
	}
	log, err := logger.New(serviceName)
	if err != nil {
		return err
	}
	shutdownObs, err := observability.Setup(ctx, serviceName, log)
	if err != nil {
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownObs(shutdownCtx); err != nil {
			log.Warn("failed to shutdown observability", zap.Error(err))
		}
	}()
	log.Info("starting service", zap.String("env", cfg.Env), zap.String("budget", cfg.ConcurrencyBudget()))
	return httpx.Run(ctx, cfg, log, func(r chi.Router) {
		register(r, cfg, log)
	})
}
