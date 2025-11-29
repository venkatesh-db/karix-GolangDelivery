# Redis 10 Lakh Data (1 Million) Real-Time Simulation in Go

This project simulates real-time data ingestion into Redis at high throughput using Go. It focuses on 10 lakh (1,000,000) live events by default using pipelining and concurrency, with configurable payload size, TTL, and batch size. Includes optional Pub/Sub events for live monitoring and a `--dry-run` guard for quick local validation.

## Prerequisites
- Go 1.22+
- A running Redis server

If you don't have Redis, run it via Docker:

```zsh
docker run -d --name redis -p 6379:6379 redis:7
```

Or start the native macOS binary (Redis 8.4.0) directly:

```zsh
# install via brew if needed
brew install redis

# start server with stock config
redis-server
```

## Quick Start

Build and run the producer with defaults (generates telemetry-style JSON payloads sized ~512 bytes):

```zsh
cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /redis10lakhsdata"
go mod tidy
# simulate 10 lakh (1,000,000) live events (tested on macOS Redis 8.4.0)
go run ./cmd/producer --count 1000000 --concurrency 12 --batch 1000 --payload-bytes 512 --ttl 0 --prefix sim --report-interval 10

# quick sanity check without touching Redis
go run ./cmd/producer --count 20000 --dry-run --report-interval 2
```

### Flags
- `--addr`: Redis address (default `127.0.0.1:6379`)
- `--password`: Redis password (empty by default)
- `--count`: Total keys to insert (default 1,000,000 = 10 lakh)
- `--concurrency`: Number of worker goroutines (default `max(CPU, 8)`)
- `--batch`: Pipeline batch size (default 1000)
- `--payload-bytes`: Payload size per key (default 512 bytes)
- `--ttl`: TTL in seconds (0 for no TTL)
- `--prefix`: Key prefix (default `sim`)
- `--publish`: Also publish events to `sim:events` channel (default false)
- `--report-interval`: Progress report interval seconds (default 5)
- `--dry-run`: Skip Redis writes for fast verification (default false)

## Live Throughput View (Optional)
Run a simple subscriber to see events:

```zsh
# Show published events count
redis-cli SUBSCRIBE sim:events
```

Or check key count:

```zsh
redis-cli INFO keyspace | grep keys
```

## Notes
- Each value is a JSON blob with event metadata (user, city, device, etc.) so downstream consumers can replay real-time signals.
- This uses Redis pipelining heavily; ensure Redis `maxclients` and networking are sufficient.
- Tune `--concurrency` and `--batch` based on your machine and Redis server.
- Use `--dry-run` when you just want to benchmark the generator without impacting Redis.
- Verified run: 1,000,000 inserts completed in ~2.0s on macOS (12 workers, batch 1000) with Redis 8.4.0 local instance.
- For very high throughput, run Redis on the same machine and avoid network bottlenecks.

## Additional Scenario: Redis vs SQLite Divergence

If you need to demonstrate cache/database inconsistency issues, check `scenarios/redis_sqlite_sync`. It replays events into Redis first and intentionally drops a subset of SQLite writes to mimic a production incident, then reports the divergent IDs.

```zsh
cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /redis10lakhsdata"
go run ./scenarios/redis_sqlite_sync --count 2000 --fail-at 150,499,1337 --random-fail-rate 0.01
```

See `scenarios/redis_sqlite_sync/README.md` for detailed instructions.

## Validation Snapshot

Both paths were executed end-to-end on macOS 14 (Redis 8.4.0):

```zsh
# 10 lakh Redis producer
go run ./cmd/producer --count 1000000 --concurrency 12 --batch 1000 --report-interval 10

# Redis vs SQLite divergence demo
go run ./scenarios/redis_sqlite_sync --count 20 --fail-at 5,12
```

The producer inserted all 1,000,000 keys in ~2 seconds, and the sync scenario reported `[evt-000005 evt-000012]` as divergent IDs, matching the injected failure points.
