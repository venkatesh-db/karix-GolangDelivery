package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/helrachar/banking/internal/api"
	appcustomer "github.com/helrachar/banking/internal/application/customer"
	"github.com/helrachar/banking/internal/config"
	"github.com/helrachar/banking/internal/infrastructure/memory"
	"github.com/helrachar/banking/internal/logging"
	"github.com/helrachar/banking/internal/server"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)

	repo := memory.NewCustomerRepository()
	service := appcustomer.NewService(repo)
	handler := api.NewCustomerHandler(service, logger)

	mux := http.NewServeMux()
	handler.Register(mux)

	decoratedHandler := server.Chain(mux,
		server.RecoveryMiddleware(logger),
		server.RequestIDMiddleware(),
		server.LoggingMiddleware(logger),
	)

	httpServer := server.NewHTTPServer(cfg.HTTP, decoratedHandler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("http server starting", "port", cfg.HTTP.Port)
		if err := httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpServer.ShutdownTimeout())
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}
