# High-Performance gRPC Ride Streaming API

This project implements a production-style bidirectional gRPC streaming API with a sharded in-memory session broker and Prometheus metrics. It also includes a matching engine stub for zone broadcasts and rider-driver matching.

## Features
- Bidirectional gRPC `RideStreamService.Connect` using envelopes (client/server).
- Sharded session broker with bounded outbound buffers and backpressure counters.
- Prometheus metrics endpoint at `:9090` (`/metrics`).
- Configurable via environment variables.

## Quick Start

1) Generate protobufs (once):
```zsh
cd "/Users/venkatesh/Golang WOW Placments/microservices/highperformancegrpcapi"
buf generate
```

2) Build and run:
```zsh
cd "/Users/venkatesh/Golang WOW Placments/microservices/highperformancegrpcapi"
GOFLAGS="" go build ./...
GRPC_LISTEN_ADDR=":7443" METRICS_LISTEN_ADDR=":9090" HTTP_LISTEN_ADDR=":8080" \
MAX_SESSIONS=1200000 OUTBOUND_BUFFER=256 HEARTBEAT_INTERVAL="5s" HEARTBEAT_TIMEOUT="15s" SHARD_COUNT=64 \
go run ./cmd/ride-stream
```

3) Check metrics:
```zsh
curl -s http://127.0.0.1:9090/metrics | head -n 50
```

## Env Vars
- `GRPC_LISTEN_ADDR` (default `:7443`): gRPC listener.
- `HTTP_LISTEN_ADDR` (default `:8080`): reserved for future HTTP gateway.
- `METRICS_LISTEN_ADDR` (default `:9090`): Prometheus metrics.
- `MAX_SESSIONS` (default `1200000`): capacity ceiling.
- `OUTBOUND_BUFFER` (default `256`): per-session outbound channel size.
- `HEARTBEAT_INTERVAL` (default `5s`): expected heartbeat cadence.
- `HEARTBEAT_TIMEOUT` (default `15s`): disconnect threshold.
- `SHARD_COUNT` (default `64`): broker shard count.

## Code Layout
- `proto/ride/v1/ride.proto` — envelope schemas and service.
- `gen/go/proto/ride/v1` — generated Go stubs.
- `internal/config` — config loader.
- `internal/telemetry` — Prometheus instruments and handler.
- `internal/stream` — `Session` and sharded `Broker`.
- `internal/matching` — matching engine stub (zone broadcast + match events).
- `internal/server` — gRPC service handler.
- `cmd/ride-stream` — service entrypoint.

## Testing the Stream
Use `grpcurl` or `evans` to open a bidirectional stream and send envelopes.

Example (evans):
```zsh
brew install evans
cd "/Users/venkatesh/Golang WOW Placments/microservices/highperformancegrpcapi"
evans --host localhost --port 7443 -r
# In evans shell:
#   package ride.v1
#   service RideStreamService
#   call Connect (bidirectional)
```

Send `ClientEnvelope` messages (heartbeat, location, status). The service will echo acks, broadcast to zone peers, and emit match events for `STATUS_LOOKING`.

## Production Notes
- Use Envoy/Linkerd for L4 consistent hashing on `user_id` across replicas.
- Enable `SO_REUSEPORT`, set TCP user timeouts, and tune node `somaxconn`.
- Scale via HPA on CPU/memory and a custom metric `ride_stream_active_sessions`.

## License
Sample code provided for reference purposes.
