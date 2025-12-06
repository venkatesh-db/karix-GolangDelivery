
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func retry(ctx context.Context, attempts int, fn func() error) error {
	backoff := 50 * time.Millisecond

	for i := 1; i <= attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err() // context canceled
		}

		err := fn()
		if err == nil {
			return nil // success
		}

		// last attempt → return error
		if i == attempts {
			return err
		}

		// exponential backoff
		select {
		case <-time.After(backoff):
			backoff = backoff * 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return errors.New("unreachable")
}

func main() {
	parent := context.Background()
	g, ctx := errgroup.WithContext(parent)

	// G1 — Debit Wallet
	g.Go(func() error {
		return retry(ctx, 3, func() error {
			fmt.Println("Debiting wallet...")
			if rand.Intn(3) == 0 {
				return nil // SUCCESS
			}
			return errors.New("wallet debit failed")
		})
	})

	// G2 — Credit Bank
	g.Go(func() error {
		return retry(ctx, 3, func() error {
			fmt.Println("Crediting bank...")
			if rand.Intn(4) == 0 {
				return nil // SUCCESS
			}
			return errors.New("bank credit failed")
		})
	})

	// G3 — Ledger Writer
	g.Go(func() error {
		return retry(ctx, 2, func() error {
			fmt.Println("Writing ledger...")
			// Ledger usually must not retry too many times
			if rand.Intn(5) == 0 {
				return nil
			}
			return errors.New("ledger write failed")
		})
	})

	// Wait for all workers
	if err := g.Wait(); err != nil {
		fmt.Println("TRANSACTION FAILED:", err)
	} else {
		fmt.Println("TRANSACTION SUCCESSFUL")
	}
}
