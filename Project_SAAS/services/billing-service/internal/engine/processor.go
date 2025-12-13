package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"project_saas/shared/pkg/concurrency"
	"project_saas/shared/pkg/config"
	"project_saas/shared/pkg/data/fake"
)

// Processor crunches large usage batches while keeping DB contention predictable.
type Processor struct {
	cfg config.ServiceConfig
	log *zap.Logger
}

// Result summarizes a billing run.
type Result struct {
	Tenant       string  `json:"tenant"`
	Processed    int64   `json:"processed"`
	OpsPerSecond float64 `json:"ops_per_sec"`
	DurationMS   int64   `json:"duration_ms"`
	Budget       string  `json:"budget"`
}

// NewProcessor builds a Processor.
func NewProcessor(cfg config.ServiceConfig, log *zap.Logger) *Processor {
	return &Processor{cfg: cfg, log: log.Named("billing-processor")}
}

// Run ingests up to total usage events (default 1M) for a tenant.
func (p *Processor) Run(ctx context.Context, tenant string, total int) (Result, error) {
	if total <= 0 {
		total = 1_000_000
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	workerLimiter := concurrency.NewLimiter(int64(max(p.cfg.MaxWorkers, 1)))
	dbLimiter := concurrency.NewLimiter(int64(max(p.cfg.MaxInFlightDBJobs, 1)))
	tracker := concurrency.NewTracker()
	start := time.Now()

	stream := fake.StreamUsage(ctx, total)
	for record := range stream {
		record := record
		workerLimiter.Go(ctx, func(ctx context.Context) error {
			if err := p.processRecord(ctx, tenant, record, dbLimiter); err != nil {
				return err
			}
			tracker.Add(1)
			return nil
		})
	}

	if err := workerLimiter.Wait(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return Result{}, concurrency.ErrExceededDeadline
		}
		return Result{}, err
	}

	processed, ops := tracker.Snapshot()
	return Result{
		Tenant:       tenant,
		Processed:    processed,
		OpsPerSecond: ops,
		DurationMS:   time.Since(start).Milliseconds(),
		Budget:       p.cfg.ConcurrencyBudget(),
	}, nil
}

func (p *Processor) processRecord(ctx context.Context, tenant string, record fake.UsageRecord, limiter *concurrency.Limiter) error {
	return limiter.Do(ctx, func(ctx context.Context) error {
		lockA, lockB := canonicalLockKeys(tenant, record.UserID)
		pace := time.Duration(100+record.Quantity%10) * time.Microsecond
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pace):
			p.log.Debug("aggregated", zap.String("lock_a", lockA), zap.String("lock_b", lockB))
			return nil
		}
	})
}

func canonicalLockKeys(tenant, user string) (string, string) {
	if tenant <= user {
		return tenant, tenant + ":" + user
	}
	return tenant + ":" + user, tenant
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var errNoRecords = errors.New("no records processed")

// Validate ensures the job actually handled data; useful in tests.
func Validate(res Result) error {
	if res.Processed == 0 {
		return fmt.Errorf("%w for tenant %s", errNoRecords, res.Tenant)
	}
	return nil
}
