package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config captures tunables for the data generator.
type Config struct {
	MongoURI         string
	Database         string
	Collection       string
	TotalRecords     int
	BatchSize        int
	Workers          int
	FailureRatio     float64
	ComplianceRatio  float64
	DuplicateRatio   float64
	StaleStatusRatio float64
	LoggerInterval   time.Duration
	ConnectTimeout   time.Duration
	OperationTimeout time.Duration
}

const (
	defaultURI         = "mongodb://localhost:27017"
	defaultDatabase    = "npcisim"
	defaultCollection  = "transactions"
	defaultRecords     = 2_500_000
	defaultBatchSize   = 2000
	defaultWorkers     = 8
	defaultFailRatio   = 0.12
	defaultCompliance  = 0.03
	defaultDuplicate   = 0.02
	defaultStaleStatus = 0.04
	defaultLogInterval = 5 * time.Second
	defaultConnTO      = 15 * time.Second
	defaultOpTO        = 30 * time.Second
)

// Load builds Config from environment variables with safe fallbacks.
func Load() (*Config, error) {
	cfg := &Config{
		MongoURI:         getString("MONGO_URI", defaultURI),
		Database:         getString("MONGO_DB", defaultDatabase),
		Collection:       getString("MONGO_COLLECTION", defaultCollection),
		TotalRecords:     getInt("TOTAL_RECORDS", defaultRecords),
		BatchSize:        getInt("BATCH_SIZE", defaultBatchSize),
		Workers:          getInt("WORKERS", defaultWorkers),
		FailureRatio:     getFloat("FAILURE_RATIO", defaultFailRatio),
		ComplianceRatio:  getFloat("COMPLIANCE_RATIO", defaultCompliance),
		DuplicateRatio:   getFloat("DUPLICATE_RATIO", defaultDuplicate),
		StaleStatusRatio: getFloat("STALE_STATUS_RATIO", defaultStaleStatus),
		LoggerInterval:   getDuration("LOGGER_INTERVAL", defaultLogInterval),
		ConnectTimeout:   getDuration("CONNECT_TIMEOUT", defaultConnTO),
		OperationTimeout: getDuration("OPERATION_TIMEOUT", defaultOpTO),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	case c.MongoURI == "":
		return fmt.Errorf("mongo uri cannot be empty")
	case c.Database == "":
		return fmt.Errorf("database name cannot be empty")
	case c.Collection == "":
		return fmt.Errorf("collection name cannot be empty")
	case c.TotalRecords <= 0:
		return fmt.Errorf("total records must be positive")
	case c.BatchSize <= 0:
		return fmt.Errorf("batch size must be positive")
	case c.Workers <= 0:
		return fmt.Errorf("workers must be positive")
	case c.FailureRatio < 0 || c.FailureRatio > 0.8:
		return fmt.Errorf("failure ratio %.2f not in range", c.FailureRatio)
	case c.ComplianceRatio < 0 || c.ComplianceRatio > 0.3:
		return fmt.Errorf("compliance ratio %.2f not in range", c.ComplianceRatio)
	case c.DuplicateRatio < 0 || c.DuplicateRatio > 0.3:
		return fmt.Errorf("duplicate ratio %.2f not in range", c.DuplicateRatio)
	case c.StaleStatusRatio < 0 || c.StaleStatusRatio > 0.5:
		return fmt.Errorf("stale status ratio %.2f not in range", c.StaleStatusRatio)
	default:
		return nil
	}
}

// RemainingRecords returns how many documents are left after `done` count.
func (c *Config) RemainingRecords(done int) int {
	left := c.TotalRecords - done
	if left < 0 {
		return 0
	}
	return left
}

func getString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func getFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed
		}
	}
	return fallback
}
