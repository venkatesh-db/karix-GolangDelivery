
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"go.uber.org/atomic"
)

//
// ---------------------------
// SHARED DATA (BANK + WALLET)
// ---------------------------
//

var walletBalance = 1000
var bankBalance = 5000

var mu sync.Mutex                 // prevents race during balance updates
var ledgerMu sync.Mutex            // prevents race during ledger writes
var processed atomic.Bool          // prevents transaction double processing
var ledger []string

//
// ---------------------------
// RETRY MECHANISM
// ---------------------------
//

func retry(ctx context.Context, attempts int, fn func() error) error {
	backoff := 50 * time.Millisecond

	for i := 1; i <= attempts; i++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := fn()
		if err == nil {
			return nil // success
		}

		if i == attempts {
			return err // final failure
		}

		// exponential backoff
		select {
		case <-time.After(backoff):
			backoff = backoff * 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

//
// ---------------------------
// BUSINESS OPERATIONS
// ---------------------------
//

func debitWallet(amount int) error {
	mu.Lock()
	defer mu.Unlock()

	if walletBalance < amount {
		return errors.New("insufficient wallet balance")
	}

	// simulate flaky API/network
	if rand.Intn(5) == 0 {
		return errors.New("wallet service timeout")
	}

	walletBalance -= amount
	return nil
}

func creditBank(amount int) error {
	mu.Lock()
	defer mu.Unlock()

	// simulate flaky CBS API
	if rand.Intn(5) != 0 {
		return errors.New("bank credit API failed")
	}

	bankBalance += amount
	return nil
}

func writeLedger(entry string) error {
	ledgerMu.Lock()
	defer ledgerMu.Unlock()

	// ledger is sensitive → only small retries allowed
	if rand.Intn(4) == 0 {
		ledger = append(ledger, entry)
		return nil
	}

	return errors.New("ledger write failed")
}

//
// ---------------------------
// TRANSACTION ORCHESTRATOR
// ---------------------------
//

func ProcessTransaction(amount int) error {

	// Prevent double processing
	if !processed.CompareAndSwap(false, true) {
		return errors.New("transaction already processed")
	}

	fmt.Println("🔵 Starting atomic transaction...")

	g, ctx := errgroup.WithContext(context.Background())

	// 1. Debit wallet
	g.Go(func() error {
		return retry(ctx, 3, func() error {
			fmt.Println("Debiting wallet...")
			return debitWallet(amount)
		})
	})

	// 2. Credit bank
	g.Go(func() error {
		return retry(ctx, 3, func() error {
			fmt.Println("Crediting bank...")
			return creditBank(amount)
		})
	})

	// 3. Ledger write
	g.Go(func() error {
		return retry(ctx, 2, func() error {
			fmt.Println("Writing ledger...")
			return writeLedger(fmt.Sprintf("₹%d transferred", amount))
		})
	})

	// Wait for all parallel tasks
	if err := g.Wait(); err != nil {
		fmt.Println("❌ Transaction FAILED:", err)

		// rollback debit if necessary
		mu.Lock()
		walletBalance += amount
		mu.Unlock()

		return err
	}

	fmt.Println("✅ Transaction SUCCESSFUL")
	return nil
}

//
// ---------------------------
// MAIN
// ---------------------------
//

func main() {
	rand.Seed(time.Now().UnixNano())

	err := ProcessTransaction(500)
	if err != nil {
		fmt.Println("Final Error:", err)
	} else {
		fmt.Println("Final Ledger:", ledger)
		fmt.Println("Wallet Balance:", walletBalance)
		fmt.Println("Bank Balance:", bankBalance)
	}
}
