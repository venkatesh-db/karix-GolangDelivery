package main

import (
	"fmt"
	"log"
)

// -------------------------
// Struct: Trade
// -------------------------
type Trade struct {
	Symbol     string
	Quantity   int
	BuyPrice   float64
	SellPrice  float64
	IsExecuted bool
}

// Method: Validate trade
func (t Trade) Validate() error {
	if t.Quantity <= 0 {
		return fmt.Errorf("invalid quantity for %s", t.Symbol)
	}
	if t.BuyPrice <= 0 || t.SellPrice <= 0 {
		return fmt.Errorf("invalid price for %s", t.Symbol)
	}
	return nil
}

// Method: Calculate profit
func (t Trade) Profit() float64 {
	return t.SellPrice - t.BuyPrice
}

// -------------------------
// Struct: Tower
// -------------------------
type Tower struct {
	Name     string
	Load     int
	Users    int
	IsActive bool
}

// Method: Tower Status
func (t Tower) Status() string {
	if t.Load < 80 {
		return "Normal"
	}
	return "Overloaded"
}

// Method: Revenue per tower
func (t Tower) Revenue(pricePerUser float64) float64 {
	return float64(t.Users) * pricePerUser
}

// -------------------------
// Main Function
// -------------------------
func main() {
	// -------------------------
	// Towers
	// -------------------------
	towers := []Tower{
		{"TowerA", 75, 120, true},
		{"TowerB", 85, 150, true},
	}

	for _, tower := range towers {
		fmt.Printf("[INFO] %s Status: %s | Revenue: ₹%.2f\n",
			tower.Name, tower.Status(), tower.Revenue(50))
	}

	// -------------------------
	// Trades
	// -------------------------
	trades := []Trade{
		{"INFY", 100, 1500, 1600, false},
		{"TCS", 50, 2000, 1950, false}, // loss trade
	}

	for _, trade := range trades {
		if err := trade.Validate(); err != nil {
			log.Printf("[ERROR] Trade validation failed: %v\n", err)
			continue
		}
		fmt.Printf("[SUCCESS] Trade %s Profit: ₹%.2f\n", trade.Symbol, trade.Profit())
	}
}
