package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config captures runtime configuration sourced from environment variables.
type Config struct {
	AppName string
	Env     string
	HTTP    HTTPConfig
}

// HTTPConfig controls HTTP server behavior.
type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Load builds a Config instance using environment variables with safe defaults.
func Load() Config {
	return Config{
		AppName: envOrDefault("APP_NAME", "customer-service"),
		Env:     envOrDefault("APP_ENV", "development"),
		HTTP: HTTPConfig{
			Port:            envOrDefault("PORT", "8080"),
			ReadTimeout:     durationFromEnv("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    durationFromEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     durationFromEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: durationFromEnv("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		d, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("invalid duration for %s, using fallback: %v", key, err)
			return fallback
		}
		return d
	}
	return fallback
}

// IntFromEnv extracts integer env vars with fallback.
func IntFromEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		log.Printf("invalid int for %s, using fallback", key)
	}
	return fallback
}
