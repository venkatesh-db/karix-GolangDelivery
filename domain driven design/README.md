# Customer Service (DDD Sample)

Domain-driven design sample for a customer onboarding service. It exposes a small REST API with in-memory storage, structured logging, and layered architecture (`domain`, `application`, `api`, `server`).

## Prerequisites
- Go 1.21+
- Free TCP port `8080` (configurable via `HTTP_PORT` env var)

## Run Tests
```bash
cd "$(dirname "$0")" && go test ./...
```

## Start the Service
```bash
cd "$(dirname "$0")" && go run ./cmd/customer-service
```
Expected startup log (fails if the port is already taken):
```
{"time":"2025-12-05T12:22:12.259469+05:30","level":"INFO","msg":"http server starting","port":"8080"}
{"time":"2025-12-05T12:22:12.260529+05:30","level":"ERROR","msg":"http server failed","error":"listen tcp :8080: bind: address already in use"}
{"time":"2025-12-05T12:22:12.260756+05:30","level":"INFO","msg":"server stopped gracefully"}
```
If you see the `address already in use` error, stop any other process bound to `:8080` or change the port:
```bash
HTTP_PORT=8081 go run ./cmd/customer-service
```

## HTTP Endpoints
| Method | Path           | Description                         |
|--------|----------------|-------------------------------------|
| GET    | `/healthz`     | Liveness check returns `{status:ok}`|
| POST   | `/customers`   | Create a customer (JSON body)       |
| GET    | `/customers`   | List all customers                  |
| GET    | `/customers/{id}` | Fetch a customer by ID          |

### Sample Requests
Create customer:
```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{"fullName":"Ada Lovelace","email":"ada@example.com","pan":"ABCDE1234F"}'
```
List customers:
```bash
curl http://localhost:8080/customers
```
Fetch by ID:
```bash
curl http://localhost:8080/customers/<customer-id>
```

All responses include a `requestId` header propagated via middleware to help trace logs.
