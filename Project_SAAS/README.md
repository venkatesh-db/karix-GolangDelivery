# Project SAAS (Go)

Modern distributed SaaS reference focused on concurrency safety, deadlock avoidance, and processing "10 lakh" (1,000,000) usage events per billing run.

## Highlights
- **Microservices**: gateway, user, subscription, billing, payment, invoicing, notification.
- **Shared tooling**: configuration loader, zap logger, concurrency limiter, fake data streaming, Postgres pool helper.
- **Subscription persistence**: subscription-service now ships with embedded migrations, repositories, and real plan activation endpoints backed by Postgres.
- **Observability**: every HTTP service is wrapped with request logging, Prometheus `/metrics`, and OpenTelemetry tracing (stdout by default or OTLP via env vars).
- **Concurrency demo**: `billing-service` streams 1M usage rows, fans out work under `MAX_WORKERS` while DB writes are throttled via `MAX_DB_JOBS`, ensuring predictable contention.
- **Persistence example**: `user-service` now provisions its own Postgres schema (embedded migrations) and exposes real CRUD endpoints for tenant users.
- **Deadlock strategy**: canonical lock ordering + advisory-limit style `Limiter.Do` around simulated DB sections.

## Quick Start
```bash
cd Project_SAAS
cp .env.example .env   # create if needed
make infra             # optional helper to run docker compose
# or manually
docker compose -f deployments/docker-compose.yml up -d

# run billing service (example)
cd services/billing-service
GOEXPERIMENT=loopvar go run ./cmd/billing-service

# run user service with Postgres persistence
cd ../user-service
HTTP_PORT=8081 go run ./cmd/user-service
```
Then trigger a billing run:
```bash
curl -X POST "http://localhost:8080/billing/tenants/acme/run?records=1000000"
```
Response includes processed count, ops/sec, and duration. Hit the user service via:
```bash
curl -X POST http://localhost:8081/tenants/acme/users \
	-H 'Content-Type: application/json' \
	-d '{"email":"alex@example.com","full_name":"Alex Admin"}'
```

## Configuration
Environment variables (per service):
- `HTTP_PORT`, `GRPC_PORT`
- `POSTGRES_URL`, `REDIS_URL`, `NATS_URL`
- `MAX_WORKERS` (goroutine fan-out), `MAX_DB_JOBS` (in-flight DB sections)
- Observability knobs: `OTEL_EXPORTER_OTLP_ENDPOINT` / `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` for remote OTLP sinks, `OBSERVABILITY_DISABLED=true` to skip tracer initialization (stdout exporter is the default otherwise).

See `docs/PROJECT_PLAN.md` and `docs/ARCHITECTURE.md` for detailed plan + diagrams.

## Next Steps
1. Implement real repositories (pgx) and transactional outbox.
2. Extend persistence patterns from `user-service`/`subscription-service` to the remaining services (billing, payments, invoicing).
3. Add ConnectRPC contracts and integrate gRPC clients.
4. Harden cross-service workflows (idempotent messaging, race/regression suites). Existing GitHub Actions CI already runs fmt/vet/tests per module.
