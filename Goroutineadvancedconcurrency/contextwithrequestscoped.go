package main

import (
	"context"
	"fmt"
)

type key string

const requestIDKey key = "reqID"

func main() {
	// Add request ID to context
	ctx := context.WithValue(context.Background(), requestIDKey, "REQ-12345")

	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	// Retrieve value
	reqID := ctx.Value(requestIDKey)
	fmt.Println("Processing with Request ID:", reqID)

	dbOperation(ctx)
}

func dbOperation(ctx context.Context) {
	reqID := ctx.Value(requestIDKey)
	fmt.Println("DB operation for Request ID:", reqID)
}
