package main

import (
	"context"
	"fmt"
	"time"
)

type key string

const requestIDKey key = "reqID"

func main() {
	// Parent context with timeout + request ID
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, requestIDKey, "REQ-99")

	process(ctx)
}

func process(ctx context.Context) {
	reqID := ctx.Value(requestIDKey)
	fmt.Println("Start process, reqID =", reqID)

	select {
	case <-time.After(3 * time.Second): // fake long work
		fmt.Println("Process finished")
	case <-ctx.Done():
		fmt.Println("❌ Cancelled:", ctx.Err(), ", reqID =", reqID)
	}
}
