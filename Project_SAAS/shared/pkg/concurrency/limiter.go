package concurrency

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Limiter bounds the number of concurrent goroutines performing heavy work (DB/IO).
type Limiter struct {
	sem      *semaphore.Weighted
	errOnce  sync.Once
	waitOnce sync.Once
	err      error
	wg       sync.WaitGroup
}

// NewLimiter creates a limiter with the provided maximum concurrency.
func NewLimiter(max int64) *Limiter {
	if max <= 0 {
		max = 1
	}
	return &Limiter{sem: semaphore.NewWeighted(max)}
}

// Go schedules fn respecting the limiter. Caller must invoke Wait after enqueuing work.
func (l *Limiter) Go(ctx context.Context, fn func(context.Context) error) {
	if fn == nil {
		return
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		if err := l.sem.Acquire(ctx, 1); err != nil {
			l.setErr(err)
			return
		}
		defer l.sem.Release(1)
		if err := fn(ctx); err != nil {
			l.setErr(err)
		}
	}()
}

// Do executes fn synchronously under the limiter budget.
func (l *Limiter) Do(ctx context.Context, fn func(context.Context) error) error {
	if fn == nil {
		return nil
	}
	if err := l.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	defer l.sem.Release(1)
	return fn(ctx)
}

// Wait blocks until all goroutines finish and returns the first error encountered.
func (l *Limiter) Wait() error {
	l.waitOnce.Do(func() {
		l.wg.Wait()
	})
	return l.err
}

func (l *Limiter) setErr(err error) {
	if err == nil {
		return
	}
	l.errOnce.Do(func() {
		l.err = err
	})
}

// ThroughputTracker estimates operations per second for large batch jobs.
type ThroughputTracker struct {
	started time.Time
	mu      sync.Mutex
	total   int64
}

// NewTracker returns a tracker.
func NewTracker() *ThroughputTracker {
	return &ThroughputTracker{started: time.Now()}
}

// Add increments processed count.
func (t *ThroughputTracker) Add(delta int64) {
	t.mu.Lock()
	t.total += delta
	t.mu.Unlock()
}

// Snapshot returns processed count and ops/sec.
func (t *ThroughputTracker) Snapshot() (count int64, opsPerSec float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	elapsed := time.Since(t.started).Seconds()
	if elapsed == 0 {
		return t.total, 0
	}
	return t.total, float64(t.total) / elapsed
}

// ErrExceededDeadline signals that work could not finish in time.
var ErrExceededDeadline = errors.New("work deadline exceeded")
