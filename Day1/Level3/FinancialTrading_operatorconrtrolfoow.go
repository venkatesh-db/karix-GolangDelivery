package main

import "log"

/*

🟢 Financial Insights

if → trading thresholds for profit/loss control.

switch → order execution mode.

for → iterate over portfolio holdings.

*/

func main() {
	currentPrice := 1575.80
	avgBuyPrice := 1500.00
	pnl := currentPrice - avgBuyPrice
	orderType := "LIMIT"

	// Conditional PnL check
	if pnl > 50 {
		log.Printf("[TRADE] Profit booking suggested | PnL=%.2f", pnl)
	} else if pnl < -20 {
		log.Printf("[TRADE] Loss beyond threshold | PnL=%.2f | Action: Hedge", pnl)
	} else {
		log.Printf("[TRADE] Holding steady | PnL=%.2f", pnl)
	}

	// Switch for order type behavior
	switch orderType {
	case "MARKET":
		log.Println("[TRADE] Executing Market Order...")
	case "LIMIT":
		log.Println("[TRADE] Executing Limit Order...")
	case "STOPLOSS":
		log.Println("[TRADE] Stop Loss Triggered...")
	default:
		log.Println("[TRADE] Unknown Order Type.")
	}

	// Loop portfolio of stocks
	symbols := []string{"INFY", "TCS", "RELIANCE", "HDFCBANK"}
	for _, s := range symbols {
		log.Printf("[TRADE] Evaluating position: %s", s)
	}
}
