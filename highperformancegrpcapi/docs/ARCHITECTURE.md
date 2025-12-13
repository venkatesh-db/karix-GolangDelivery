# High-Performance Ride Streaming API

## Objectives
- Real-time bidirectional streaming between riders, drivers, and dispatch services using gRPC.
- WebSocket fan-out for browser / edge clients while keeping gRPC as the system-of-record transport.
- Sustain 1M concurrent connections with predictable tail latency under p95 < 150 ms.

## Component Overview
1. **RideStream gRPC Service**
   - Bidirectional `Connect` RPC exchanging `ClientEnvelope` and `ServerEnvelope` messages.
   - Implements session auth, flow control, and backpressure using bounded channels.
   - Uses tuned `grpc.Server` options (buffer sizes, keepalive, connection limits).
2. **Session Broker**
   - Lock-free `sync.Map` index keyed by user or ride identifiers.
   - Per-session `ringbuffer` backed by pooled byte slices to minimize allocations.
   - Supports multicast topics (city/zone) and direct unicast for acknowledgements.
3. **WebSocket Gateway**
   - HTTP server terminating TLS and upgrading to `nhooyr.io/websocket` connections.
   - Bridges WebSocket frames to the broker via the same envelope format (JSON).
   - Implements adaptive batching for outbound messages to align with websocket flush cadence.
4. **Telemetry + Safeguards**
   - Prometheus metrics for connection counts, queue depth, stream latency.
   - Zap structured logging with per-request correlation IDs.
   - Feature flags for sampling, synthetic load, and chaos testing hooks.

## Scaling Strategy
- $N_{conn}$ split across shards by hashing `user_id` mod shard-count; each shard owns its own broker instance.
- Envoy/Linkerd performs L4 load balancing with consistent hashing on `user_id` to preserve affinity.
- Use `SO_REUSEPORT` + multiple listener goroutines pinned via `automaxprocs` for CPU saturation control.
- Apply CELL architecture (region -> zone -> pod) to isolate failure domains.

## Data Contracts
```
rpc Connect(stream ClientEnvelope) returns (stream ServerEnvelope);
```
- `ClientEnvelope`: location updates, ride status, telematics heartbeats.
- `ServerEnvelope`: matchmaking decisions, ETA refresh, acknowledgements, broadcast events.
- Every envelope carries `corr_id` and lamport clock for total ordering when merged downstream.

## Reliability Features
- Session level heartbeats every 5s; disconnect after 3 misses.
- Backpressure signals propagate by pausing reads on both gRPC and WebSocket sides when outbound buffers exceed 75%.
- Pluggable persistence sink (Kafka / Pulsar) for guaranteed delivery; default noop stub in sample code.

## Deployment Notes
- Run gRPC server on :7443 (h2c behind ingress) and WebSocket gateway on :8080 (HTTP/1.1).
- HPA scales pods based on CPU (60%), memory (70%), and custom metric `active_sessions`.
- Enable `SO_KEEPALIVE`, set TCP user timeouts, and fine-tune `net.core.somaxconn` at the node layer.
