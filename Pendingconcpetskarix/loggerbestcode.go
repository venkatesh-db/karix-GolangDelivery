package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var baseLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

type ctxKey string

const loggerKey ctxKey = "logger"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		reqLogger := baseLogger.With(
			"request_id", uuid.NewString(),
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
		)

		ctx := context.WithValue(r.Context(), loggerKey, reqLogger)

		reqLogger.Info("Incoming Request")

		next.ServeHTTP(w, r.WithContext(ctx))

		reqLogger.Info("Completed Request")
	})
}

func GetLogger(r *http.Request) *slog.Logger {
	logger := r.Context().Value(loggerKey)
	if logger != nil {
		return logger.(*slog.Logger)
	}
	return baseLogger
}

func HelloAPI(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r)
	logger.Info("Inside HelloAPI handler", "action", "processing")

	w.Write([]byte("Hello from API"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", LoggingMiddleware(http.HandlerFunc(HelloAPI)))

	fmt.Println("smiling jamon faces 😊 Server started at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		baseLogger.Error("Server stopped unexpectedly", "error", err.Error())
	}
}
