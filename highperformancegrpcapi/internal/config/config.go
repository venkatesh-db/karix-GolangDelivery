package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config captures the tunable settings for the streaming stack.
type Config struct {
	GRPCListenAddr     string
	HTTPListenAddr     string
	MetricsListenAddr  string
	MaxSessions        int
	OutboundBufferSize int
	HeartbeatInterval  time.Duration
	HeartbeatTimeout   time.Duration
	ShardCount         int
}

// Load parses the process environment and returns a populated Config.
func Load() (Config, error) {
	cfg := Config{
		GRPCListenAddr:     valueOrDefault("GRPC_LISTEN_ADDR", ":7443"),
		HTTPListenAddr:     valueOrDefault("HTTP_LISTEN_ADDR", ":8080"),
		MetricsListenAddr:  valueOrDefault("METRICS_LISTEN_ADDR", ":9090"),
		MaxSessions:        intFromEnv("MAX_SESSIONS", 1200000),
		OutboundBufferSize: intFromEnv("OUTBOUND_BUFFER", 256),
		HeartbeatInterval:  durationFromEnv("HEARTBEAT_INTERVAL", 5*time.Second),
		HeartbeatTimeout:   durationFromEnv("HEARTBEAT_TIMEOUT", 15*time.Second),
		ShardCount:         intFromEnv("SHARD_COUNT", 64),
	}

	if cfg.OutboundBufferSize < 32 {
		return Config{}, fmt.Errorf("OUTBOUND_BUFFER must be >= 32")
	}
	if cfg.ShardCount <= 0 {
		return Config{}, fmt.Errorf("SHARD_COUNT must be positive")
	}
	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func intFromEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func durationFromEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
//