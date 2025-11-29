package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	redis "github.com/redis/go-redis/v9"
)

func main() {
	var (
		addr       = flag.String("addr", "127.0.0.1:6379", "Redis address")
		password   = flag.String("password", "", "Redis password")
		sqlitePath = flag.String("sqlite", "./temp/sync_issue.db", "SQLite database path")
		prefix     = flag.String("prefix", "sync", "Redis key prefix")
		count      = flag.Int("count", 2000, "Number of events to emit")
		failAt     = flag.String("fail-at", "150,499,1337", "Comma-separated event indexes that will skip SQLite writes")
		randomFail = flag.Float64("random-fail-rate", 0.0, "Probability (0-1) of random SQLite write skip per event")
		reset      = flag.Bool("reset", true, "Reset Redis keys and SQLite file before running")
	)
	flag.Parse()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: *addr, Password: *password})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("unable to reach redis: %v", err)
	}

	if *reset {
		if err := os.Remove(*sqlitePath); err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to remove sqlite file: %v", err)
		}
		if err := clearRedisPrefix(ctx, rdb, *prefix); err != nil {
			log.Fatalf("failed to clear redis prefix: %v", err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(*sqlitePath), 0o755); err != nil {
		log.Fatalf("failed to create sqlite dir: %v", err)
	}

	db, err := sql.Open("sqlite3", *sqlitePath)
	if err != nil {
		log.Fatalf("failed to open sqlite: %v", err)
	}
	defer db.Close()

	if err := bootstrapSQLite(db); err != nil {
		log.Fatalf("failed to bootstrap sqlite: %v", err)
	}

	injectedFailures := parseFailList(*failAt)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Printf("Starting sync simulator: events=%d redis=%s sqlite=%s failPoints=%v randomFailRate=%.2f", *count, *addr, *sqlitePath, injectedFailures, *randomFail)

	expectedIDs := make([]string, 0, *count)
	redisWrites := 0
	sqliteWrites := 0

	for i := 0; i < *count; i++ {
		id := fmt.Sprintf("evt-%06d", i)
		payload := fmt.Sprintf(`{"id":"%s","total":%d,"ts":%d}`, id, rnd.Intn(10_000), time.Now().UnixMilli())
		key := fmt.Sprintf("%s:%s", *prefix, id)

		if err := rdb.Set(ctx, key, payload, 0).Err(); err != nil {
			log.Fatalf("redis write failed at %s: %v", id, err)
		}
		redisWrites++
		expectedIDs = append(expectedIDs, id)

		if shouldSkipSQLite(i, injectedFailures, *randomFail, rnd) {
			log.Printf("simulated production fault: redis has %s but sqlite skipped", id)
			continue
		}

		if err := insertEvent(db, id, payload); err != nil {
			log.Fatalf("sqlite insert failed at %s: %v", id, err)
		}
		sqliteWrites++
	}

	summary := compareStores(ctx, rdb, db, *prefix, expectedIDs)
	log.Printf("Run summary: redis=%d sqlite=%d divergent=%d", redisWrites, sqliteWrites, len(summary.DivergentIDs))
	if len(summary.DivergentIDs) > 0 {
		log.Printf("Sample divergent IDs (first 10): %v", summary.DivergentIDs[:min(10, len(summary.DivergentIDs))])
	} else {
		log.Printf("No divergence detected")
	}
}

func bootstrapSQLite(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		payload TEXT NOT NULL,
		inserted_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func insertEvent(db *sql.DB, id, payload string) error {
	_, err := db.Exec(`INSERT INTO events(id, payload) VALUES(?, ?)`, id, payload)
	return err
}

func clearRedisPrefix(ctx context.Context, rdb *redis.Client, prefix string) error {
	var cursor uint64
	pattern := fmt.Sprintf("%s:*", prefix)
	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 500).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func parseFailList(raw string) map[int]struct{} {
	result := make(map[int]struct{})
	if strings.TrimSpace(raw) == "" {
		return result
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var idx int
		if _, err := fmt.Sscanf(part, "%d", &idx); err == nil {
			result[idx] = struct{}{}
		}
	}
	return result
}

func shouldSkipSQLite(i int, failList map[int]struct{}, randomRate float64, rnd *rand.Rand) bool {
	if _, ok := failList[i]; ok {
		return true
	}
	if randomRate <= 0 {
		return false
	}
	return rnd.Float64() < randomRate
}

type syncSummary struct {
	DivergentIDs []string
}

func compareStores(ctx context.Context, rdb *redis.Client, db *sql.DB, prefix string, expected []string) syncSummary {
	sqliteIDs, err := loadSQLiteIDs(db)
	if err != nil {
		log.Fatalf("failed to load sqlite ids: %v", err)
	}
	divergent := make([]string, 0)
	for _, id := range expected {
		if _, ok := sqliteIDs[id]; !ok {
			key := fmt.Sprintf("%s:%s", prefix, id)
			exists, err := rdb.Exists(ctx, key).Result()
			if err != nil {
				log.Fatalf("redis exists failed: %v", err)
			}
			if exists > 0 {
				divergent = append(divergent, id)
			}
		}
	}
	return syncSummary{DivergentIDs: divergent}
}

func loadSQLiteIDs(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query(`SELECT id FROM events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make(map[string]struct{})
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = struct{}{}
	}
	return ids, rows.Err()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
