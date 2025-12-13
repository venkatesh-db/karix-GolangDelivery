package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/venkatesh/order-service/internal/api"
	"github.com/venkatesh/order-service/internal/app"
	"github.com/venkatesh/order-service/internal/eventbus"
	"github.com/venkatesh/order-service/internal/infrastructure"
	"github.com/venkatesh/order-service/internal/readmodel"
)

func main() {
	store := infrastructure.NewInMemoryStore()
	bus := eventbus.New(4)
	projection := readmodel.NewOrdersProjection()
	bus.Subscribe(eventbus.WildcardEvent, projection.Handle)

	service := app.NewOrderService(store, bus)
	httpServer := api.NewServer(service, projection)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(httpServer.Router()),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("order-service listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server gracefully stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s completed in %s", r.Method, r.URL.Path, time.Since(start))
	})
}
