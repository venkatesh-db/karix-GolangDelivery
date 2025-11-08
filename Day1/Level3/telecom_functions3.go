package main

import (
	"fmt"
	"log"
)

// -------------------------
// 1️⃣ Multiple return - Telecom Tower Status
// -------------------------
func towerStatus(tower string, load int) (string, bool) {
	if load < 80 {
		return "Normal", true
	}
	return "Overloaded", false
}

// -------------------------
// 2️⃣ Variadic - Total connections across towers
// -------------------------
func totalConnections(connections ...int) int {
	total := 0
	for _, c := range connections {
		total += c
	}
	return total
}

// -------------------------
// 3️⃣ Inline Function - Revenue Calculation per Tower
// -------------------------
var calcRevenue = func(users int, pricePerUser float64) float64 {
	return float64(users) * pricePerUser
}

// -------------------------
// 4️⃣ Recursive - Total users in tower batches
// -------------------------
func totalUsersBatch(usersPerTower []int, index int) int {
	if index < 0 {
		return 0
	}
	return usersPerTower[index] + totalUsersBatch(usersPerTower, index-1)
}

// -------------------------
// 5️⃣ Defer - Tower Monitoring Cleanup
// -------------------------
func monitorTower(tower string) {
	defer fmt.Println("[INFO] Finished monitoring:", tower)
	fmt.Println("[INFO] Monitoring started:", tower)
}

// -------------------------
// 6️⃣ Panic & Recover - Validate Tower Data
// -------------------------
func processTower(tower string, load int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Tower %s processing failed: %v\n", tower, r)
		}
	}()
	if load < 0 || load > 100 {
		panic("Invalid tower load value")
	}
	status, active := towerStatus(tower, load)
	fmt.Printf("[SUCCESS] Tower %s Status: %s | Active: %t\n", tower, status, active)
}

func main() {
	// Tower loads
	towers := map[string]int{
		"TowerA": 75,
		"TowerB": 85,
		"TowerC": 60,
		"TowerD": -5, // Invalid load to trigger panic/recover
	}

	// 1️⃣ Multiple Return + Panic/Recover
	for tower, load := range towers {
		processTower(tower, load)
	}

	// 2️⃣ Variadic - Total Connections
	connections := []int{120, 90, 150}
	fmt.Println("[INFO] Total connections across towers:", totalConnections(connections...))

	// 3️⃣ Inline Function - Revenue per tower
	usersPerTower := []int{120, 90, 150}
	pricePerUser := 50.0
	for i, users := range usersPerTower {
		revenue := calcRevenue(users, pricePerUser)
		fmt.Printf("[INFO] Revenue for Tower %d: ₹%.2f\n", i+1, revenue)
	}

	// 4️⃣ Recursive - Total Users Batch
	totalUsers := totalUsersBatch(usersPerTower, len(usersPerTower)-1)
	fmt.Println("[INFO] Total users across all towers (recursive):", totalUsers)

	// 5️⃣ Defer - Monitor Towers
	for tower := range towers {
		monitorTower(tower)
	}
}

