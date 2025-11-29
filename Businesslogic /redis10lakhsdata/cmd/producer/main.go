package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	mathrand "math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var (
	eventTypes = []string{"signup", "purchase", "page_view", "heartbeat", "logout"}
	cities     = []string{"Mumbai", "Delhi", "Hyderabad", "Chennai", "Bengaluru", "Pune", "Kolkata"}
	devices    = []string{"ios", "android", "web", "tablet", "smart-tv"}
)

type telemetryEvent struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Event     string  `json:"event"`
	Amount    float64 `json:"amount"`
	City      string  `json:"city"`
	Device    string  `json:"device"`
	Metadata  string  `json:"metadata"`
	Padding   string  `json:"padding,omitempty"`
	Timestamp int64   `json:"ts"`
}

func main() {
	var (
		addr           = flag.String("addr", "127.0.0.1:6379", "Redis address")
		password       = flag.String("password", "", "Redis password")
		count          = flag.Int64("count", 1_000_000, "Total keys to insert (10 lakh = 1,000,000)")
		concurrency    = flag.Int("concurrency", max(int64(runtime.NumCPU()), 8), "Number of concurrent workers")
		batchSize      = flag.Int("batch", 1000, "Pipeline batch size")
		payloadBytes   = flag.Int("payload-bytes", 512, "Target payload size per key in bytes")
		ttlSeconds     = flag.Int("ttl", 0, "TTL for each key in seconds (0 = no TTL)")
		prefix         = flag.String("prefix", "sim", "Key prefix")
		publishEvents  = flag.Bool("publish", false, "Publish events to <prefix>:events channel")
		reportInterval = flag.Int("report-interval", 5, "Progress report interval in seconds")
		dryRun         = flag.Bool("dry-run", false, "Simulate without touching Redis (for quick verification)")
	)
	flag.Parse()

	log.Printf("Redis real-time simulator starting: addr=%s count=%d concurrency=%d batch=%d payload=%d ttl=%d prefix=%s publish=%v dry=%v",
		*addr, *count, *concurrency, *batchSize, *payloadBytes, *ttlSeconds, *prefix, *publishEvents, *dryRun)

	ctx := context.Background()
	var rdb *redis.Client
	if !*dryRun {
		rdb = redis.NewClient(&redis.Options{Addr: *addr, Password: *password})
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("redis ping failed: %v", err)
		}
	}

	var inserted int64
	var counter int64
	start := time.Now()

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(*reportInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				curr := atomic.LoadInt64(&inserted)
				elapsed := time.Since(start).Seconds()
				throughput := float64(curr) / math.Max(elapsed, 0.001)
				log.Printf("Progress: inserted=%d elapsed=%.1fs throughput=%.0f ops/s", curr, elapsed, throughput)
			case <-done:
				return
			}
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(*concurrency)

	for workerID := 0; workerID < *concurrency; workerID++ {
		go func(id int) {
			defer wg.Done()
			rnd := mathrand.New(mathrand.NewSource(time.Now().UnixNano() + int64(id)))
			for {
				startIdx := atomic.AddInt64(&counter, int64(*batchSize)) - int64(*batchSize)
				if startIdx >= *count {
					return
				}
				endIdx := startIdx + int64(*batchSize)
				if endIdx > *count {
					endIdx = *count
				}

				if *dryRun {
					atomic.AddInt64(&inserted, endIdx-startIdx)
					continue
				}

				pipe := rdb.Pipeline()
				ttl := time.Duration(*ttlSeconds) * time.Second

				for i := startIdx; i < endIdx; i++ {
					key := fmt.Sprintf("%s:%d", *prefix, i)
					payload := buildEventPayload(rnd, i, *payloadBytes)
					if ttl > 0 {
						pipe.SetEx(ctx, key, payload, ttl)
					} else {
						pipe.Set(ctx, key, payload, 0)
					}
					if *publishEvents {
						pipe.Publish(ctx, *prefix+":events", fmt.Sprintf("insert:%d", i))
					}
				}

				if _, err := pipe.Exec(ctx); err != nil {
					log.Printf("worker %d pipeline error: %v", id, err)
					time.Sleep(200 * time.Millisecond)
					if _, err = pipe.Exec(ctx); err != nil {
						log.Printf("worker %d failed after retry: %v", id, err)
						return
					}
				}

				atomic.AddInt64(&inserted, endIdx-startIdx)
			}
		}(workerID)
	}

	wg.Wait()
	close(done)

	elapsed := time.Since(start).Seconds()
	log.Printf("Completed: inserted=%d in %.1fs (%.0f ops/s)", inserted, elapsed, float64(inserted)/math.Max(elapsed, 0.001))
}

func buildEventPayload(r *mathrand.Rand, seq int64, targetBytes int) []byte {
	evt := telemetryEvent{
		ID:        fmt.Sprintf("evt-%d", seq),
		UserID:    fmt.Sprintf("user-%d", seq%500_000),
		Event:     eventTypes[r.Intn(len(eventTypes))],
		Amount:    math.Round((r.Float64()*9999+1)*100) / 100,
		City:      cities[r.Intn(len(cities))],
		Device:    devices[r.Intn(len(devices))],
		Metadata:  randomAlphaNum(r, 32),
		Timestamp: time.Now().UTC().UnixMilli(),
	}
	payload, _ := json.Marshal(evt)
	if len(payload) >= targetBytes {
		return payload
	}
	evt.Padding = randomAlphaNum(r, (targetBytes-len(payload))*2)
	payload, _ = json.Marshal(evt)
	return payload
}

func randomAlphaNum(r *mathrand.Rand, length int) string {
	if length <= 0 {
		return ""
	}
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	b.Grow(length)
	for i := 0; i < length; i++ {
		b.WriteByte(alphabet[r.Intn(len(alphabet))])
	}
	return b.String()
}

func max(a, b int64) int {
	if a > b {
		return int(a)
	}
	return int(b)
}
