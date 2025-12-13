# Stage 1: Build
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o order-service

# Stage 2: Runtime
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 8082
ENTRYPOINT ["/app/order-service"]