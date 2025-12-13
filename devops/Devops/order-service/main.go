package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type orderResponse struct {
	ID    int    `json:"id"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

var (
	httpRequestsTotal   = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func main() {
	logger := log.New(os.Stdout, "order-service: ", log.LstdFlags|log.Lshortfile)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderResponse{ID: 1, Item: "Laptop", Price: 1200})
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		logger.Printf("method=%s path=%s status=%d latency_ms=%.2f", r.Method, r.URL.Path, 200, time.Since(start).Seconds()*1000)
	})

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":" + port}

	go func() {
		logger.Printf("Starting server on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}