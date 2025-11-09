
package main



/*

✅ Industry-Grade Quality Features

| Feature         | What’s Improved                 | Why It’s Production-Ready                 |
| --------------- | ------------------------------- | ----------------------------------------- |
| Structs         | `Trade` struct                  | Clean, encapsulated trade data            |
| Logging         | `log.Printf` + `[INFO]/[ERROR]` | Realistic logging for production          |
| Error Handling  | Multiple return + panic/recover | Prevent crashes & handle invalid trades   |
| Variadic        | Total Volume                    | Flexible dynamic input                    |
| Inline Function | Profit calculation              | Reusable, clean logic                     |
| Recursive       | Compounded profit               | Demonstrates repeated computations        |
| Defer           | `processOrder`                  | Ensures cleanup/logging at end of process |
| Looping         | For all trades                  | Batch processing for microservices        |


*/

import (
	"errors"
	"fmt"
	"log"
)

// -------------------------
// Trade Struct for Industry-Grade Model
// -------------------------
type Trade struct {
	Symbol   string
	Quantity int
	BuyPrice float64
	SellPrice float64
}

// -------------------------
// 1️⃣ Multiple Return - Validate Trade
// -------------------------
func validateTrade(trade Trade) (string, error) {
	if trade.Quantity <= 0 {
		return "", errors.New("quantity must be greater than zero")
	}
	if trade.BuyPrice <= 0 || trade.SellPrice <= 0 {
		return "", errors.New("price must be greater than zero")
	}
	return "Trade valid", nil
}

// -------------------------
// 2️⃣ Variadic Function - Total Volume
// -------------------------
func totalVolume(volumes ...int) int {
	total := 0
	for _, v := range volumes {
		total += v
	}
	return total
}

// -------------------------
// 3️⃣ Inline Function - Profit Calculation
// -------------------------
var calcProfit = func(buy, sell float64) float64 {
	return sell - buy
}

// -------------------------
// 4️⃣ Recursive Function - Compounded Profit
// -------------------------
func compoundedProfit(profit float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	return profit + compoundedProfit(profit, n-1)
}

// -------------------------
// 5️⃣ Defer Statement - Cleanup Logs
// -------------------------
func processOrder(trade Trade) {
	defer fmt.Println("[INFO] Finished processing order:", trade.Symbol)
	fmt.Println("[INFO] Processing order for:", trade.Symbol)
}

// -------------------------
// 6️⃣ Panic & Recover - Execute Trade Safely
// -------------------------
func executeTrade(trade Trade) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Trade %s failed: %v\n", trade.Symbol, r)
		}
	}()

	if trade.Quantity <= 0 || trade.BuyPrice <= 0 || trade.SellPrice <= 0 {
		panic("Invalid trade parameters")
	}

	profit := calcProfit(trade.BuyPrice, trade.SellPrice)
	fmt.Printf("[SUCCESS] Trade %s executed successfully | Profit: ₹%.2f\n", trade.Symbol, profit)
}

// -------------------------
// Main Function - Industry Simulation
// -------------------------
func main() {
	trades := []Trade{
		{"INFY", 100, 1500, 1600},
		{"TCS", 0, 2000, 2100},        // Invalid quantity
		{"RELI", 50, 2500, 2450},      // Loss trade
		{"WIPRO", 30, 1000, 1200},
	}

	// Validate Trades
	for _, trade := range trades {
		msg, err := validateTrade(trade)
		if err != nil {
			log.Printf("[ERROR] Validation failed for %s: %v\n", trade.Symbol, err)
		} else {
			fmt.Println("[INFO]", msg, "for", trade.Symbol)
		}
	}

	// Total Volume Calculation (Variadic)
	volumes := []int{100, 50, 30, 200}
	fmt.Println("[INFO] Total traded volume:", totalVolume(volumes...))

	// Process Orders with Defer
	for _, trade := range trades {
		processOrder(trade)
	}

	// Execute Trades Safely (Panic & Recover)
	for _, trade := range trades {
		executeTrade(trade)
	}

	// Recursive Compounded Profit Example
	profit := compoundedProfit(500, 3)
	fmt.Printf("[INFO] Compounded Profit for 3 trades: ₹%.2f\n", profit)
}


