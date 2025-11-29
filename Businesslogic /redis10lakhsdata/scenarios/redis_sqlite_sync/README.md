# Redis vs SQLite Sync Divergence Simulator

This scenario reproduces a common production issue where Redis accepts writes but SQLite (or any relational store) silently misses a subset of events due to transaction failures or crash loops. The program inserts events into Redis first and conditionally skips the SQLite insert to mimic an outage, then reports the divergent keys.

## Features
- Deterministic failure points (`--fail-at`) to recreate exact incidents.
- Optional random failure rate to simulate flaky DB connections.
- Automatic cleanup (`--reset`) removes previous keys under the prefix and deletes the SQLite file.
- Summary report listing IDs present in Redis but absent in SQLite.

## Run It
```zsh
cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /redis10lakhsdata"
go run ./scenarios/redis_sqlite_sync --count 2000 --fail-at 150,499,1337 --random-fail-rate 0.01 --prefix sync
```

Sample output:
```
2025/11/29 14:05:42 simulated production fault: redis has evt-000150 but sqlite skipped
...
2025/11/29 14:05:45 Run summary: redis=2000 sqlite=1967 divergent=33
2025/11/29 14:05:45 Sample divergent IDs (first 10): [evt-000150 evt-000499 evt-001337 ...]
```

## Flags
- `--addr`: Redis address (default `127.0.0.1:6379`).
- `--password`: Redis password.
- `--sqlite`: Path to the SQLite DB file (default `./temp/sync_issue.db`).
- `--prefix`: Redis key prefix (default `sync`).
- `--count`: Number of events to emit (default 2000).
- `--fail-at`: Comma-separated event indexes that will skip SQLite writes.
- `--random-fail-rate`: Probability (0-1) that an event randomly skips SQLite.
- `--reset`: Drop existing Redis keys/SQLite file before running (default true).

## Cleanup
Set `--reset=false` if you want to inspect the leftover SQLite DB or Redis keys between runs. Keys live under `<prefix>:evt-*` and SQLite rows are in the `events` table.
