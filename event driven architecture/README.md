# Order Service (Event-Driven Sample)

This repository contains a production-style Go service that models the lifecycle of e-commerce orders using Domain-Driven Design, Event Sourcing, and an in-memory event bus.

## Architecture Overview

- **Domain layer (`internal/domain`)**: Order aggregate, line items, and strongly typed domain events.
- **Application layer (`internal/app`)**: Command DTOs and `OrderService` orchestrating event-store persistence plus publishing.
- **Infrastructure (`internal/infrastructure`)**: In-memory event store implementing optimistic concurrency.
- **Event bus (`internal/eventbus`)**: Lightweight async pub/sub dispatcher (swap for Kafka/NATS in real deployments).
- **Read model (`internal/readmodel`)**: Projection that builds a query-optimized view for HTTP responses.
- **HTTP API (`internal/api`)**: REST endpoints for commands and querying the projection.

Everything wires together in `cmd/orderservice/main.go`, exposing port `8080`.

## Prerequisites

- Go 1.22+
- Port `8080` available on your machine

## Running the Service

```bash
cd "/Users/venkatesh/Golang WOW Placments/microservices/event driven architecture"
go run ./cmd/orderservice
```

You should see:

```
2025/12/05 12:17:47 order-service listening on :8080
```

Press `Ctrl+C` to stop.

## Sample Workflow (cURL)

1. **Health check**
   ```bash
   curl -s http://localhost:8080/health
   # {"status":"ok"}
   ```
2. **Place an order**
   ```bash
   curl -s -X POST http://localhost:8080/orders \
     -H 'Content-Type: application/json' \
     -d '{"order_id":"ord-001","customer_id":"cust-123","items":[{"SKU":"sku-1","Quantity":2,"UnitPriceCents":1500}]}'
   # {"order_id":"ord-001"}
   ```
3. **Authorize payment**
   ```bash
   curl -s -X POST http://localhost:8080/orders/ord-001/payment \
     -H 'Content-Type: application/json' \
     -d '{"payment_id":"pay-001","amount_cents":3000}'
   # {"status":"accepted"}
   ```
4. **Reserve inventory**
   ```bash
   curl -s -X POST http://localhost:8080/orders/ord-001/reserve \
     -H 'Content-Type: application/json' \
     -d '{"reservation_id":"res-001"}'
   # {"status":"accepted"}
   ```
5. **Ship order**
   ```bash
   curl -s -X POST http://localhost:8080/orders/ord-001/ship \
     -H 'Content-Type: application/json' \
     -d '{"tracking_number":"1Z999","carrier":"UPS"}'
   # {"status":"accepted"}
   ```
6. **Query read model**
   ```bash
   curl -s http://localhost:8080/orders
   # [{"OrderID":"ord-001","CustomerID":"cust-123","Status":"shipped","TotalCents":3000,"LastUpdated":1764917328}]

   curl -s http://localhost:8080/orders/ord-001
   # {"OrderID":"ord-001","CustomerID":"cust-123","Status":"shipped","TotalCents":3000,"LastUpdated":1764917328}
   ```

## Development Tips

- Format code with `gofmt`: `find . -name '*.go' -print0 | xargs -0 gofmt -w`
- Run compilation/tests: `go build ./...` or `go test ./...`
- Swap the in-memory store/event bus with durable infrastructure (e.g., Postgres + Kafka) for persistence and out-of-process consumers.
