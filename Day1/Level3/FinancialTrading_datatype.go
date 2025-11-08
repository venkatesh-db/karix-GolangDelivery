/*

🟢 Financial Production Context:

float64 ensures price precision (₹1875.60)

int32 controls memory footprint for quantity

Logging structured for downstream ELK/Splunk parsers

*/

package main

import (
	"log"
)

func main() {
	// 💰 Financial trading production variables
	var (
		orderID        string  = "ORDX1245"
		tradePrice     float64 = 1875.60
		quantity       int32   = 100
		isOrderFilled  bool    = true
		exchangeSymbol string  = "NSE:INFY"
	)

	// Structured order event log
	log.Printf("[TRADE] OrderID=%s | Symbol=%s | Price=%.2f | Qty=%d | Filled=%t",
		orderID, exchangeSymbol, tradePrice, quantity, isOrderFilled)
}
