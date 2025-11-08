package main

import "fmt"

/*

| Feature         | Use in Trading     | Explanation                                       |
| --------------- | ------------------ | ------------------------------------------------- |
| Multiple Return | Trade Validation   | Returns message + validity                        |
| Variadic        | Total Volume       | Sum across multiple trades dynamically            |
| Inline Function | Profit Calculation | Quick reusable profit logic                       |
| Recursive       | Compounded Profit  | Calculate repeated profits across multiple trades |
| Defer           | Order Cleanup      | Ensure logs or DB closure after order             |
| Panic & Recover | Invalid Trades     | Protect system from crashing on bad input         |

*/

// 1️⃣ Multiple return values - Trade validation
func validateTrade(symbol string, quantity int, price float64) (string, bool) {
	if quantity <= 0 || price <= 0 {
		return "Invalid trade parameters", false
	}
	return "Trade valid", true
}

// 2️⃣ Variadic function - Total traded volume
func totalVolume(volumes ...int) int {
	total := 0
	for _, v := range volumes {
		total += v
	}
	return total
}

// 3️⃣ Inline / anonymous function - Profit calculation
var calcProfit = func(buy, sell float64) float64 {
	return sell - buy
}

// 4️⃣ Recursive function - Compute compounded profit (simplified)
func compoundedProfit(profit float64, n int) float64 {
	if n == 0 {
		return 0
	}
	return profit + compoundedProfit(profit, n-1)
}

// 5️⃣ Defer - Cleanup logs after trade processing
func processOrder(symbol string) {
	defer fmt.Println("Finished processing order:", symbol)
	fmt.Println("Processing order for:", symbol)
}

// 6️⃣ Panic & Recover - Handle invalid trades
func executeTrade(symbol string, amount float64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Trade %s failed: %v\n", symbol, r)
		}
	}()
	if amount <= 0 {
		panic("Trade amount must be greater than zero")
	}
	fmt.Printf("Trade %s executed successfully: ₹%.2f\n", symbol, amount)
}

func main() {
	// -------------------------
	// Multiple return example
	status, ok := validateTrade("INFY", 100, 1500)
	fmt.Println("Validation:", status, "✅", ok)

	status, ok = validateTrade("TCS", 0, 2000)
	fmt.Println("Validation:", status, "✅", ok)

	// -------------------------
	// Variadic function example
	vol := totalVolume(100, 200, 50, 150)
	fmt.Println("Total Traded Volume:", vol)

	// -------------------------
	// Inline function example
	profit := calcProfit(1000, 1200)
	fmt.Println("Profit from trade:", profit)

	// -------------------------
	// Recursive function example
	compProfit := compoundedProfit(500, 3)
	fmt.Println("Compounded Profit over 3 trades:", compProfit)

	// -------------------------
	// Defer statement example
	processOrder("RELI")

	// -------------------------
	// Panic & Recover example
	executeTrade("INFY", 5000)
	executeTrade("TCS", -2000) // triggers panic & recover
}
