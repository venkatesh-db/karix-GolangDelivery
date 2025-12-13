package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once
var loadErr error

// ServiceConfig captures the knobs every service cares about.
type ServiceConfig struct {
	ServiceName       string
	HTTPPort          string
	GRPCPort          string
	PostgresURL       string
	RedisURL          string
	NATSURL           string
	AuthSecret        string
	MaxWorkers        int
	MaxInFlightDBJobs int
	Env               string
}

// Load reads environment variables (optionally from .env) once per process.
func Load(service string) (ServiceConfig, error) {
	loadOnce.Do(func() {
		if _, err := os.Stat(".env"); err == nil {
			loadErr = godotenv.Load()
		}
	})
	if loadErr != nil {
		return ServiceConfig{}, loadErr
	}

	cfg := ServiceConfig{
		ServiceName:       service,
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		GRPCPort:          getEnv("GRPC_PORT", "9090"),
		PostgresURL:       getEnv("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/saas?sslmode=disable"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/0"),
		NATSURL:           getEnv("NATS_URL", "nats://localhost:4222"),
		AuthSecret:        getEnv("AUTH_SECRET", "insecure-dev-secret"),
		Env:               getEnv("APP_ENV", "development"),
		MaxWorkers:        getEnvInt("MAX_WORKERS", 32),
		MaxInFlightDBJobs: getEnvInt("MAX_DB_JOBS", 8),
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func getEnvInt(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return def
}

// ConcurrencyBudget returns a formatted string used for logging.
func (c ServiceConfig) ConcurrencyBudget() string {
	return fmt.Sprintf("workers=%d db_jobs=%d", c.MaxWorkers, c.MaxInFlightDBJobs)
}
