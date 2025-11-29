package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/glebarez/sqlite"
)

type stepTimer struct{ started time.Time }

func (s *stepTimer) start()              { s.started = time.Now() }
func (s *stepTimer) stop() time.Duration { return time.Since(s.started) }

func main() {
	var (
		mode    = flag.String("mode", "opt", "Mode: 'opt' (optimized) or 'unopt' (unoptimized)")
		rows    = flag.Int("rows", 1000000, "Number of rows to insert (10 lakhs = 1,000,000)")
		dbPath  = flag.String("db", "./sqlite_10lakhs.db", "SQLite database file path")
		batch   = flag.Int("batch", 10000, "Batch size for optimized inserts (ignored in unoptimized mode)")
		fresh   = flag.Bool("fresh", true, "If true, removes existing DB file before running")
		vacuum  = flag.Bool("vacuum", false, "Run VACUUM at the end (optimized mode)")
		verbose = flag.Bool("v", false, "Verbose output (e.g., explain plan)")
	)
	flag.Parse()

	if *fresh {
		_ = os.Remove(*dbPath)
	}

	if err := os.MkdirAll(filepath.Dir(*dbPath), 0o755); err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	dsn := sqliteDSN(*dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// conservative connection limits
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()

	summary := make(map[string]time.Duration)

	s := &stepTimer{}
	var runErr error

	s.start()
	if strings.EqualFold(*mode, "unopt") {
		runErr = runUnoptimized(ctx, db, *rows, *verbose)
	} else if strings.EqualFold(*mode, "opt") {
		runErr = runOptimized(ctx, db, *rows, *batch, *vacuum, *verbose)
	} else {
		runErr = fmt.Errorf("unknown mode: %s", *mode)
	}
	summary["total"] = s.stop()

	if runErr != nil {
		log.Fatalf("run error: %v", runErr)
	}

	printSummary(*mode, *rows, summary)
}

func sqliteDSN(path string) string {
	// modernc.org/sqlite uses driver name "sqlite" and supports file: URIs.
	// shared cache enables WAL better, but not required.
	return fmt.Sprintf("file:%s?cache=shared", path)
}

func mustExec(ctx context.Context, db *sql.DB, q string) {
	if _, err := db.ExecContext(ctx, q); err != nil {
		log.Fatalf("exec failed: %s: %v", q, err)
	}
}

func runUnoptimized(ctx context.Context, db *sql.DB, n int, verbose bool) error {
	// Default PRAGMAs (do nothing). Demonstrate anti-patterns: per-row Exec, no transaction, index before load.
	s := &stepTimer{}

	// Schema + index BEFORE load (unoptimized)
	s.start()
	mustExec(ctx, db, `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at INTEGER NOT NULL
	)`)
	mustExec(ctx, db, `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`) // costly before bulk load
	schemaDur := s.stop()

	// Inserts: one Exec per row, implicit transaction per statement
	s.start()
	for i := 1; i <= n; i++ {
		name := fmt.Sprintf("User %d", i)
		email := fmt.Sprintf("user%07d@example.com", i)
		created := time.Now().Unix() + int64(i%86400)
		if _, err := db.ExecContext(ctx,
			"INSERT INTO users (id, name, email, created_at) VALUES (?, ?, ?, ?)",
			i, name, email, created,
		); err != nil {
			return fmt.Errorf("insert row %d: %w", i, err)
		}
		if i%100000 == 0 {
			log.Printf("unoptimized: inserted %d rows...", i)
		}
	}
	insertDur := s.stop()

	// Example queries
	qDurCount, qDurPoint, qDurLike, err := runQueries(ctx, db, n, verbose)
	if err != nil {
		return err
	}

	printSection("UNOPTIMIZED", map[string]time.Duration{
		"schema+preindex": schemaDur,
		"insert":          insertDur,
		"query_count":     qDurCount,
		"query_point":     qDurPoint,
		"query_like":      qDurLike,
	})
	return nil
}

func runOptimized(ctx context.Context, db *sql.DB, n, batch int, doVacuum, verbose bool) error {
	// Apply performance PRAGMAs suitable for one-time bulk load.
	applyOptimPragmas(ctx, db)

	s := &stepTimer{}

	// Schema without indexes first
	s.start()
	mustExec(ctx, db, `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		created_at INTEGER NOT NULL
	)`)
	schemaDur := s.stop()

	// Bulk inserts within a single transaction + prepared statement
	if batch <= 0 {
		batch = 10000
	}

	s.start()
	if err := bulkInsert(ctx, db, n, batch); err != nil {
		return err
	}
	insertDur := s.stop()

	// Create indexes AFTER load
	s.start()
	mustExec(ctx, db, `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`) // cheaper post-load
	indexDur := s.stop()

	// Optional VACUUM to apply page_size changes and compact file (costly time-wise)
	var vacuumDur time.Duration
	if doVacuum {
		s.start()
		mustExec(ctx, db, `VACUUM`)
		vacuumDur = s.stop()
	}

	// Example queries
	qDurCount, qDurPoint, qDurLike, err := runQueries(ctx, db, n, verbose)
	if err != nil {
		return err
	}

	printSection("OPTIMIZED", map[string]time.Duration{
		"schema":      schemaDur,
		"insert_bulk": insertDur,
		"index_post":  indexDur,
		"vacuum":      vacuumDur,
		"query_count": qDurCount,
		"query_point": qDurPoint,
		"query_like":  qDurLike,
	})
	return nil
}

func applyOptimPragmas(ctx context.Context, db *sql.DB) {
	// Aggressive PRAGMAs for one-time bulk load. If durability matters during load,
	// avoid synchronous=OFF and consider NORMAL or default.
	stmts := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=OFF",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA locking_mode=EXCLUSIVE",
		"PRAGMA mmap_size=268435456", // 256MB
		"PRAGMA cache_size=-200000",  // ~200MB cache (KB when negative)
		"PRAGMA page_size=4096",
		"PRAGMA foreign_keys=OFF",
	}
	for _, q := range stmts {
		mustExec(ctx, db, q)
	}
}

func bulkInsert(ctx context.Context, db *sql.DB, n, batch int) error {
	// Single transaction; batched loops; prepared statement for speed
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users (id, name, email, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	start := 1
	for start <= n {
		end := start + batch - 1
		if end > n {
			end = n
		}
		for i := start; i <= end; i++ {
			name := fmt.Sprintf("User %d", i)
			email := fmt.Sprintf("user%07d@example.com", i)
			created := time.Now().Unix() + int64(i%86400)
			if _, err := stmt.ExecContext(ctx, i, name, email, created); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("insert row %d: %w", i, err)
			}
		}
		if (end)%100000 == 0 || end == n {
			log.Printf("optimized: inserted %d rows...", end)
		}
		start = end + 1
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func runQueries(ctx context.Context, db *sql.DB, n int, verbose bool) (time.Duration, time.Duration, time.Duration, error) {
	s := &stepTimer{}

	// COUNT(*)
	s.start()
	var c int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&c); err != nil {
		return 0, 0, 0, err
	}
	countDur := s.stop()
	if c != n {
		return 0, 0, 0, fmt.Errorf("count mismatch: got %d want %d", c, n)
	}

	// Point lookup by indexed column
	needle := fmt.Sprintf("user%07d@example.com", 7*n/10) // 70th percentile
	s.start()
	var id int
	var name string
	var email string
	var created int64
	if err := db.QueryRowContext(ctx, "SELECT id, name, email, created_at FROM users WHERE email=?", needle).Scan(&id, &name, &email, &created); err != nil {
		return 0, 0, 0, err
	}
	pointDur := s.stop()

	// LIKE scan (non-indexed pattern) - demonstrates slower path
	s.start()
	likePrefix := "User 9" // matches many rows when n is large
	rows, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM users WHERE name LIKE ?", likePrefix+"%")
	if err != nil {
		return 0, 0, 0, err
	}
	var likeCount int
	if rows.Next() {
		if err := rows.Scan(&likeCount); err != nil {
			rows.Close()
			return 0, 0, 0, err
		}
	}
	rows.Close()
	likeDur := s.stop()

	if verbose {
		printExplain(ctx, db, "SELECT id FROM users WHERE email=?", needle)
		printExplain(ctx, db, "SELECT COUNT(*) FROM users WHERE name LIKE ?", likePrefix+"%")
	}

	return countDur, pointDur, likeDur, nil
}

func printExplain(ctx context.Context, db *sql.DB, q string, arg any) {
	log.Printf("EXPLAIN QUERY PLAN for: %s", q)
	rows, err := db.QueryContext(ctx, "EXPLAIN QUERY PLAN "+q, arg)
	if err != nil {
		log.Printf("  explain error: %v", err)
		return
	}
	defer rows.Close()
	var id, parent, notused int
	var detail string
	for rows.Next() {
		if err := rows.Scan(&id, &parent, &notused, &detail); err != nil {
			log.Printf("  explain scan error: %v", err)
			return
		}
		log.Printf("  %d %d %d: %s", id, parent, notused, detail)
	}
}

func printSection(title string, parts map[string]time.Duration) {
	log.Printf("==== %s ====", title)
	keys := []string{"schema+preindex", "schema", "insert", "insert_bulk", "index_post", "vacuum", "query_count", "query_point", "query_like"}
	for _, k := range keys {
		if d, ok := parts[k]; ok && d > 0 {
			log.Printf("%-16s: %s", k, d)
		}
	}
}

func printSummary(mode string, n int, total map[string]time.Duration) {
	log.Printf("==== SUMMARY ====")
	log.Printf("mode: %s, rows: %d, total: %s", mode, n, total["total"])
}
