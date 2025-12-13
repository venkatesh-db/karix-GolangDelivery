package logging

import (
	"log/slog"
	"os"
	"strings"
)

// New constructs a structured logger with environment-aware log levels.
func New(env string) *slog.Logger {
	level := slog.LevelInfo
	if strings.EqualFold(env, "development") {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
