package main

/*

Function → encapsulates business logic

Returns (result, error) → standard Go idiom

Handles success/failure cleanly

Easy to extend for microservices (e.g., batch transaction processing)

*/


import (
	"errors"
	"fmt"
)

// Function to process payment
func processTransaction(txnID string, amount float64, success bool) (string, error) {
	if !success {
		return "", errors.New("transaction failed due to insufficient balance")
	}
	return fmt.Sprintf("Transaction %s of ₹%.2f completed", txnID, amount), nil
}

func main() {
	transactions := []struct {
		ID      string
		Amount  float64
		Success bool
	}{
		{"TXN101", 5000, true},
		{"TXN102", 12000, false},
		{"TXN103", 7500, true},
	}

	for _, txn := range transactions {
		message, err := processTransaction(txn.ID, txn.Amount, txn.Success)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", txn.ID, err)
		} else {
			fmt.Printf("✅ %s\n", message)
		}
	}
}


