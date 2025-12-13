# Go SaaS Platform Plan

## Goals
- Scaffold a modern distributed SaaS reference in Go that highlights concurrency control, deadlock avoidance, and database-safe access patterns when processing up to 1,000,000 ("10 lakh") records.
- Provide clear service boundaries for User, Subscription, Billing, Payment, Invoicing, and Gateway components.
- Showcase communication mix: REST for external APIs, gRPC for internal RPCs, and NATS (or similar) for async flows.
- Deliver runnable examples plus documentation so engineers can extend quickly.

## Repository Layout
```
Project_SAAS/
  go.work                # Go workspace aggregating all services
  shared/                # Shared libraries (config, logging, db, concurrency helpers)
  services/
    gateway/
    user-service/
    subscription-service/
    billing-service/
    payment-service/
    invoicing-service/
    notification-service/
  deployments/
    docker-compose.yml   # Local stack (Postgres, Redis, NATS)
  docs/
    PROJECT_PLAN.md
    ARCHITECTURE.md
```

## Key Technical Decisions
- **Language**: Go 1.22 with `toolchain go1.22` to enable MVS across modules.
- **Dependencies**: Chi (REST), ConnectRPC (gRPC+HTTP bridge), NATS JetStream client, pgx for Postgres, Redis go client, OpenTelemetry SDK.
- **Concurrency primitives**: `errgroup`, `semaphore`, custom worker pools, idempotency tokens stored in Redis, transactional outbox pattern.
- **Data simulation**: `shared/data/fake` will stream 1,000,000 usage rows via bounded channels; processors use configurable goroutine limits and report throughput metrics.
- **Testing**: Provide load-test harness using `go test ./... -run TestThroughput` and Vegeta scripts under `tests/load`.
- **Configuration**: `shared/config` loads from env + `.env`. Each service has `cmd/<service>/main.go`.

## Next Steps
1. Initialize Go workspace and module stubs per service.
2. Add shared concurrency toolkit + fake database generator.
3. Implement representative handlers demonstrating deadlock-safe usage of Postgres (advisory locks + SKIP LOCKED) and concurrency limiting.
4. Document run instructions and case study insights.
