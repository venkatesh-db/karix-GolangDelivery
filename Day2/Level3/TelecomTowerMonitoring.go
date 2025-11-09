package main

import "fmt"

// Struct for Telecom Tower
type Tower struct {
	Name     string
	Load     int // percentage load
	Users    int
	Revenue  float64
	IsActive bool
}

// Method: Check tower status
func (t Tower) Status() string {
	if t.Load < 80 {
		return "Normal"
	}
	return "Overloaded"
}

// Method: Calculate revenue per tower
func (t Tower) RevenuePerUser(price float64) float64 {
	return float64(t.Users) * price
}

func main() {
	// Initializing struct with values
	towerA := Tower{
		Name:     "TowerA",
		Load:     75,
		Users:    120,
		Revenue:  0,
		IsActive: true,
	}

	fmt.Printf("[INFO] %s Status: %s\n", towerA.Name, towerA.Status())
	fmt.Printf("[INFO] Revenue for %s: ₹%.2f\n", towerA.Name, towerA.RevenuePerUser(50))
}
