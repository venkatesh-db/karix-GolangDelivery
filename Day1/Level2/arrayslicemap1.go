
package main

import "fmt"

func main() {
	// =========================
	// 1️⃣ Goibibo / IRCTC Booking System
	// =========================
	fmt.Println("=== Goibibo / IRCTC Booking System ===")

	// Array: Fixed trains for a route
	trains := [3]string{"Train101", "Train102", "Train103"}
	fmt.Println("Available Trains:", trains)

	// Slice: Dynamic bookings
	bookings := []string{"PNR123", "PNR124"}
	bookings = append(bookings, "PNR125")
	fmt.Println("Active Bookings:", bookings)

	// Map: PNR → Seat Number
	seatMap := map[string]int{
		"PNR123": 12,
		"PNR124": 18,
	}
	seatMap["PNR125"] = 25

	fmt.Println("PNR to Seat Map:")
	for pnr, seat := range seatMap {
		fmt.Printf("  %s -> Seat %d\n", pnr, seat)
	}

	// Check a PNR
	checkPNR := "PNR124"
	if seat, ok := seatMap[checkPNR]; ok {
		fmt.Printf("%s found with Seat %d ✅\n", checkPNR, seat)
	} else {
		fmt.Printf("%s not found ❌\n", checkPNR)
	}

	// =========================
	// 2️⃣ Google Pay Transaction System
	// =========================
	fmt.Println("\n=== Google Pay Transaction System ===")

	// Array: Fixed supported banks
	banks := [3]string{"SBI", "HDFC", "ICICI"}
	fmt.Println("Supported Banks:", banks)

	// Slice: Dynamic daily transactions
	transactions := []string{"TXN101", "TXN102"}
	transactions = append(transactions, "TXN103")
	fmt.Println("Daily Transactions:", transactions)

	// Map: TransactionID → Status
	transactionStatus := map[string]string{
		"TXN101": "Success",
		"TXN102": "Pending",
	}
	transactionStatus["TXN103"] = "Failed"

	fmt.Println("Transaction Status Map:")
	for txn, status := range transactionStatus {
		fmt.Printf("  %s -> %s\n", txn, status)
	}

	// Check a transaction
	checkTxn := "TXN102"
	if status, ok := transactionStatus[checkTxn]; ok {
		fmt.Printf("%s Status: %s ✅\n", checkTxn, status)
	} else {
		fmt.Printf("%s not found ❌\n", checkTxn)
	}
}
