# SQLite 10 Lakhs (1,000,000) Records — Go: Unoptimized vs Optimized

This small Go program demonstrates inserting and querying ~10 lakhs (1,000,000) rows in SQLite using two approaches:
- Unoptimized: per-row inserts without transactions and with indexes created before loading.
- Optimized: bulk insert inside one transaction using a prepared statement, fast PRAGMAs, and indexes created after loading.

It uses the pure-Go driver `modernc.org/sqlite` (no CGO needed).

## Quick Start

```bash
# macOS zsh
cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /sqlite10lakhsdata"

# (Optional) initialize go modules if needed
# go mod tidy

# Run optimized path for 1,000,000 rows
go run . -mode opt -rows 1000000 -fresh=true -v

# Run unoptimized path for 1,000,000 rows
go run . -mode unopt -rows 1000000 -fresh=true -v
```

Flags:
- `-mode`: `opt` (optimized) or `unopt` (unoptimized)
- `-rows`: number of rows to insert (set to 1,000,000 for 10 lakhs)
- `-fresh`: if `true`, remove existing DB file before running
- `-batch`: batch size for optimized inserts (default 10,000)
- `-vacuum`: if `true` (optimized), runs `VACUUM` at the end to compact
- `-db`: database file path (default `./sqlite_10lakhs.db`)
- `-v`: verbose (prints EXPLAIN QUERY PLAN for sample queries)

## What the Program Does

- Creates a `users` table with columns `(id INTEGER PRIMARY KEY, name TEXT, email TEXT, created_at INTEGER)`.
- Inserts `N` rows with synthetic data: unique email per row and a simple created_at.
- Runs a few queries and times them:
  - `SELECT COUNT(*) FROM users`
  - `SELECT id,name,email,created_at FROM users WHERE email=?` (indexed)
  - `SELECT COUNT(*) FROM users WHERE name LIKE 'User 9%'` (non-indexed pattern)

## Unoptimized Path (what NOT to do for bulk load)
- Creates index before loading data (`CREATE INDEX ... ON users(email)`).
- Inserts each row with `db.Exec` without an explicit transaction (implicit transaction per statement).
- Leaves default PRAGMAs in place.

This causes:
- Massive transaction overhead (each row is its own transaction).
- Index maintenance cost for every row, significantly slowing inserts.

## Optimized Path (what to do)
- Sets PRAGMAs suited for bulk load:
  - `journal_mode=WAL`
  - `synchronous=OFF` (only if acceptable during load; improves speed)
  - `temp_store=MEMORY`, `locking_mode=EXCLUSIVE`
  - `mmap_size=256MB`, `cache_size≈200MB`, `page_size=4096`
  - `foreign_keys=OFF` (no FKs here)
- Creates the table first, but defers indexes until after the bulk load.
- Performs all inserts inside one transaction with a prepared statement.
- Optionally runs `VACUUM` at the end to compact/align page size.

This reduces:
- Write amplification and per-row transaction cost.
- Index maintenance overhead during insert phase.

## Expected Differences
While exact timings depend on hardware and SSD, the **optimized** path typically achieves:
- 10x–100x faster insert throughput vs unoptimized.
- Much smaller file growth rate during bulk insert.
- Faster indexed lookup (`WHERE email=?`) after index is created.

## Notes and Safety
- `synchronous=OFF` improves insert speed at the cost of durability during the load. If you cannot tolerate data loss on a crash during bulk load, use `synchronous=NORMAL` or keep defaults and accept slower speed.
- For production, turn durability-related PRAGMAs back to safer values after loading.
- If you change `page_size`, run `VACUUM` once to apply it.

## Example Commands for 10 Lakhs

```bash
# Optimized (recommended for bulk load)
go run . -mode opt -rows 1000000 -fresh=true -batch 20000 -v

# Unoptimized (for comparison)
go run . -mode unopt -rows 1000000 -fresh=true -v
```

## Troubleshooting
- If you see module download errors, ensure your internet access is available and run `go mod tidy`.
- If you prefer CGO-based driver, swap imports to `github.com/mattn/go-sqlite3` and open with `sql.Open("sqlite3", "./sqlite_10lakhs.db")`. Ensure you have a working C toolchain for CGO.
