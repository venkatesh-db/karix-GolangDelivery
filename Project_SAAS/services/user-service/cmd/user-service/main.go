package main

import (
	"context"
	"os/signal"
	"syscall"

	"project_saas/shared/pkg/bootstrap"

	"project_saas/services/user-service/internal/http/routes"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := bootstrap.RunHTTPService(ctx, "user-service", routes.Register); err != nil {
		panic(err)
	}
}
