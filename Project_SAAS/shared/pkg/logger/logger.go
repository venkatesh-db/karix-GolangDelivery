package logger

import (
	"os"

	"go.uber.org/zap"
)

// New creates a structured zap logger with sane defaults for local development.
func New(service string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	if os.Getenv("APP_ENV") == "development" {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.InitialFields = map[string]interface{}{
		"service": service,
	}
	return cfg.Build()
}
