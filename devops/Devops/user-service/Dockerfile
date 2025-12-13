# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o user-service

# Stage 2: Runtime
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 8081
ENTRYPOINT ["/app/user-service"]