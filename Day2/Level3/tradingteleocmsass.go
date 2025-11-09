/*

Multi-Industry Production Features

| Module           | Key Struct     | Features                                  |
| ---------------- | -------------- | ----------------------------------------- |
| Trading          | `Trade`        | Validate, Profit, Panic/Recover           |
| Telecom          | `TelecomTower` | Status, Revenue                           |
| Automotive + IoT | `Vehicle`      | Fuel Status, Telemetry                    |
| Cloud/SaaS       | `CloudServer`  | Health, Alerts                            |
| Common           | Functions      | Variadic, Recursive, Inline functions     |
| Production       | Logging        | `log.Printf` for errors/alerts            |
| Cleanup          | Defer          | Applied in panic/recover (Trading module) |


*/

package main

import (
	"fmt"
	"log"
)

// ==============================
// Trading Module
// ==============================
type Trade struct {
	Symbol     string
	Quantity   int
	BuyPrice   float64
	SellPrice  float64
	IsExecuted bool
}

// Validate trade data
func (t Trade) Validate() error {
	if t.Quantity <= 0 {
		return fmt.Errorf("invalid quantity for %s", t.Symbol)
	}
	if t.BuyPrice <= 0 || t.SellPrice <= 0 {
		return fmt.Errorf("invalid price for %s", t.Symbol)
	}
	return nil
}

// Calculate profit for trade
func (t Trade) Profit() float64 {
	return t.SellPrice - t.BuyPrice
}

// Execute trade safely with panic/recover
func executeTrade(trade Trade) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[TRADING ERROR] Trade %s failed: %v", trade.Symbol, r)
		}
	}()
	if err := trade.Validate(); err != nil {
		panic(err)
	}
	fmt.Printf("[TRADING SUCCESS] %s executed. Profit: ₹%.2f\n", trade.Symbol, trade.Profit())
}

// ==============================
// Telecom Module
// ==============================
type TelecomTower struct {
	Name     string
	Load     int
	Users    int
	IsActive bool
}

// Check tower status
func (t TelecomTower) Status() string {
	if t.Load < 80 {
		return "Normal"
	}
	return "Overloaded"
}

// Revenue calculation
func (t TelecomTower) Revenue(pricePerUser float64) float64 {
	return float64(t.Users) * pricePerUser
}

// ==============================
// Automotive + Mobility + IoT Module
// ==============================
type Vehicle struct {
	LicensePlate string
	FuelLevel    float64 // percentage
	Location     string
	IsActive     bool
}

// Check vehicle readiness
func (v Vehicle) Status() string {
	if v.FuelLevel < 15 {
		return "Low Fuel"
	}
	return "Ready"
}

// IoT Telemetry report
func (v Vehicle) Telemetry() {
	fmt.Printf("[VEHICLE REPORT] %s | Fuel: %.1f%% | Location: %s | Status: %s\n",
		v.LicensePlate, v.FuelLevel, v.Location, v.Status())
}

// ==============================
// Cloud + SaaS Module
// ==============================
type CloudServer struct {
	Name        string
	CPUUsage    int
	MemoryUsage int
	IsActive    bool
}

// Check server health
func (s CloudServer) Health() string {
	if s.CPUUsage < 75 && s.MemoryUsage < 80 {
		return "Healthy"
	}
	return "Critical"
}

// Alert server issues
func (s CloudServer) Alert() {
	if s.Health() == "Critical" {
		log.Printf("[CLOUD ALERT] %s Overloaded! CPU: %d%%, Memory: %d%%\n", s.Name, s.CPUUsage, s.MemoryUsage)
	}
}

// ==============================
// Main Function - Industry Simulation
// ==============================
func main() {
	fmt.Println("=========== TRADING MODULE ===========")
	trades := []Trade{
		{"INFY", 100, 1500, 1600, false},
		{"TCS", 50, 2000, 1950, false}, // loss trade
	}
	for _, trade := range trades {
		executeTrade(trade)
	}

	fmt.Println("\n=========== TELECOM MODULE ===========")
	towers := []TelecomTower{
		{"TowerA", 75, 120, true},
		{"TowerB", 85, 150, true},
	}
	for _, tower := range towers {
		fmt.Printf("[TELECOM INFO] %s | Status: %s | Revenue: ₹%.2f\n",
			tower.Name, tower.Status(), tower.Revenue(50))
	}

	fmt.Println("\n=========== AUTOMOTIVE / IoT MODULE ===========")
	vehicles := []Vehicle{
		{"KA-01-AB-1234", 80, "Bangalore", true},
		{"KA-02-CD-5678", 10, "Chennai", true}, // low fuel
	}
	for _, vehicle := range vehicles {
		vehicle.Telemetry()
	}

	fmt.Println("\n=========== CLOUD / SAAS MODULE ===========")
	servers := []CloudServer{
		{"Server-A", 65, 70, true},
		{"Server-B", 85, 90, true}, // critical
	}
	for _, server := range servers {
		fmt.Printf("[CLOUD INFO] %s | Health: %s\n", server.Name, server.Health())
		server.Alert()
	}

	fmt.Println("\n=========== BATCH PROCESSING EXAMPLES ===========")
	// Recursive: Total users across towers
	usersPerTower := []int{120, 150}
	totalUsers := recursiveTotalUsers(usersPerTower, len(usersPerTower)-1)
	fmt.Println("[BATCH INFO] Total users across all towers (recursive):", totalUsers)

	// Variadic: Total revenue across towers
	totalRevenue := totalRevenuePerTower(50, 120, 150)
	fmt.Println("[BATCH INFO] Total revenue across towers (variadic): ₹", totalRevenue)

	// Inline function: quick profit for a single trade
	calcSingleProfit := func(buy, sell float64) float64 {
		return sell - buy
	}
	fmt.Println("[INLINE FUNCTION] Quick profit:", calcSingleProfit(1000, 1200))
}

// Recursive function for total users
func recursiveTotalUsers(users []int, index int) int {
	if index < 0 {
		return 0
	}
	return users[index] + recursiveTotalUsers(users, index-1)
}

// Variadic function for total revenue
func totalRevenuePerTower(pricePerUser float64, users ...int) float64 {
	total := 0.0
	for _, u := range users {
		total += float64(u) * pricePerUser
	}
	return total
}
