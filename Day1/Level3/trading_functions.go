

package main

import (
	"errors"
	"fmt"
)

// Function: process trade with alert
func processTrade(symbol string, profit float64) (string, error) {
	if profit < 0 {
		return "", errors.New("loss alert: monitor position")
	}
	return fmt.Sprintf("Trade %s successful ✅ Profit: %.2f", symbol, profit), nil
}

func main() {
	trades := map[string]float64{
		"INFY":  5000,
		"TCS":  -2000,
		"RELI": 7000,
		"WIPRO": -1500,
	}

	for symbol, profit := range trades {
		msg, err := processTrade(symbol, profit)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", symbol, err)
		} else {
			fmt.Println(msg)
		}
	}
}

