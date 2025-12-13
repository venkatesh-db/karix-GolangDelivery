# Go SaaS Microservices Architecture

## Core Services
| Service | Scope | Tech Notes |
| --- | --- | --- |
| Gateway | Public REST edge, rate limiting, JWT validation | Chi router, forwards to internal services via service mesh |
| User Service | Tenants, org users, RBAC, SCIM sync | Exposes `/tenants/{id}/users`, stores data in Postgres with embedded migrations |
| Subscription Service | Plans, metering, entitlement checks | Keeps seat counts, metered usage windows |
| Billing Service | Aggregates up to 1M usage records per run, orchestrates invoices | Demonstrates concurrency limiting + deadlock avoidance |
| Payment Service | Payment intents, PSP callbacks, dunning | Integrates with Stripe/Adyen via outbound workers |
| Invoicing Service | Generates PDF invoices, tax rules | Emits invoice events and stores PDFs in blob |
| Notification Service | Sends e-mail/SMS/webhooks | Subscribes to billing + invoicing events |

## Communication
- **North-south**: REST/JSON APIs via Gateway.
- **East-west**: Internal gRPC (ConnectRPC) plus NATS JetStream for async sagas (not fully wired yet, but placeholders exist).
- **Events**: Billing runs emit `invoice.ready` events consumed by notification + payment watchers.

## Concurrency + Deadlock Controls
- Each service has `MAX_WORKERS` and `MAX_DB_JOBS` env knobs, surfaced through `config.ServiceConfig`.
- Shared `concurrency.Limiter` offers `Go` (async pool) + `Do` (sync section) around `semaphore.Weighted`.
- Billing processor streams synthetic `10,00,000` usage records using bounded channel to stress goroutine scheduling.
- DB contention avoided via `canonicalLockKeys(tenant,user)` so locks acquired in deterministic order; combined with limited `MaxInFlightDBJobs`, this mirrors Postgres advisory locks or Cockroach key ordering.
- `ThroughputTracker` surfaces ops/sec for monitoring dashboards.

## Local Infrastructure
`deployments/docker-compose.yml` provides Postgres, Redis, NATS JetStream, and CockroachDB single-node. All services default to `localhost` endpoints.

## Extending the Skeleton
1. Duplicate service template (module, go.mod, cmd/main).
2. Wire gRPC handlers via ConnectRPC (pending stub).
3. Replace fake data generator with actual repositories built atop `shared/pkg/postgres` and `shared/pkg/redis`.
4. Add saga coordinator using NATS subjects for `billing.run.requested`, `payment.intent.succeeded`, etc.
5. Expand `tests/` with load + race-condition suites (Vegeta/k6, `go test -race`).

## Case Study: Concurrency Regression & Fix
- **Problem**: During a mock "ConcurFlow" tenant migration, importing 10 lakh expense rows caused billing workers to saturate Postgres, leading to deadlocks between tenant summary rows and payment intents.
- **Detection**: Traces showed random lock ordering on `(tenant_id, user_id)` keys and unlimited goroutines hammering the DB pool.
- **Remediation baked into this repo**:
	- `concurrency.Limiter` enforces goroutine + DB budgets per service.
	- `canonicalLockKeys` ensures deterministic lock acquisition order.
	- Billing run is capped by context deadline (2 minutes) returning `ErrExceededDeadline` if breached, preventing cascading failures.
	- Result payload exposes throughput so autoscalers can react before saturation.
