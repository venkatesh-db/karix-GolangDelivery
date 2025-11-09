
package main

import (
	"fmt"
	"log"
	"time"
)

// ---------------------------
// Common Interfaces
// ---------------------------

// Reporter — entities that can produce a report snapshot and an ID.
type Reporter interface {
	ID() string
	Report() (string, error)
}

// Executor — entities that can Execute an operation and return result or error.
type Executor interface {
	Exec() (string, error)
}

// HealthChecker — generic health check
type HealthChecker interface {
	Check() error
}

// ---------------------------
// 1) Trading — Level 3 Industry
// ---------------------------
type Trade struct {
	Symbol   string
	Quantity int
	BuyPrice float64
	SellPrice float64
}

func (t Trade) ID() string { return "trade:" + t.Symbol }

func (t Trade) Report() (string, error) {
	if t.Quantity <= 0 {
		return "", fmt.Errorf("invalid quantity for %s", t.Symbol)
	}
	return fmt.Sprintf("Trade %s qty=%d buy=%.2f sell=%.2f", t.Symbol, t.Quantity, t.BuyPrice, t.SellPrice), nil
}

func (t Trade) Exec() (string, error) {
	// basic validation
	if t.BuyPrice <= 0 || t.SellPrice <= 0 {
		return "", fmt.Errorf("invalid prices for %s", t.Symbol)
	}
	profit := t.SellPrice - t.BuyPrice
	return fmt.Sprintf("Executed %s profit=₹%.2f", t.Symbol, profit), nil
}

// ---------------------------
// 2) Telecom — Level 3 Industry
// ---------------------------
type Tower struct {
	IDField string
	Load    int // percent
	Users   int
}

func (tw Tower) ID() string { return "tower:" + tw.IDField }

func (tw Tower) Report() (string, error) {
	if tw.Users < 0 {
		return "", fmt.Errorf("negative user count for %s", tw.IDField)
	}
	status := "Normal"
	if tw.Load >= 85 {
		status = "Degraded"
	}
	return fmt.Sprintf("Tower %s load=%d users=%d status=%s", tw.IDField, tw.Load, tw.Users, status), nil
}

func (tw Tower) Check() error {
	if tw.Load >= 95 {
		return fmt.Errorf("critical load for %s", tw.IDField)
	}
	return nil
}

// ---------------------------
// 3) SaaS Product — Level 3 Industry
// ---------------------------
type SaaSApp struct {
	Name       string
	ActiveSubs int
	CPU        int // %
	Mem        int // %
}

func (s SaaSApp) ID() string { return "saas:" + s.Name }

func (s SaaSApp) Report() (string, error) {
	if s.ActiveSubs < 0 {
		return "", fmt.Errorf("invalid active subs for %s", s.Name)
	}
	health := "Healthy"
	if s.CPU > 80 || s.Mem > 85 {
		health = "Unhealthy"
	}
	return fmt.Sprintf("App %s subs=%d cpu=%d mem=%d health=%s", s.Name, s.ActiveSubs, s.CPU, s.Mem, health), nil
}

func (s SaaSApp) Check() error {
	if s.CPU > 95 {
		return fmt.Errorf("cpu overload for %s", s.Name)
	}
	return nil
}

// ---------------------------
// 4) Supply Chain — Level 3 Industry
// ---------------------------
type SupplyNode struct {
	NodeID     string
	Inventory  int
	Throughput int // units per hour
}

func (sn SupplyNode) ID() string { return "supply:" + sn.NodeID }

func (sn SupplyNode) Report() (string, error) {
	if sn.Inventory < 0 {
		return "", fmt.Errorf("invalid inventory at %s", sn.NodeID)
	}
	risk := "OK"
	if sn.Inventory < 50 {
		risk = "LowStock"
	}
	return fmt.Sprintf("Node %s inv=%d throughput=%d risk=%s", sn.NodeID, sn.Inventory, sn.Throughput, risk), nil
}

func (sn SupplyNode) Exec() (string, error) {
	// example: trigger replenishment if inventory low
	if sn.Inventory < 20 {
		return "", fmt.Errorf("replenish request sent for %s", sn.NodeID)
	}
	return fmt.Sprintf("Supply %s healthy", sn.NodeID), nil
}

// ---------------------------
// 5) Volkswagen Automotive — Level 3 Industry (connected car example)
// ---------------------------
type Vehicle struct {
	VIN       string
	FuelPct   float64
	Connected bool
	Location  string
}

func (v Vehicle) ID() string { return "vehicle:" + v.VIN }

func (v Vehicle) Report() (string, error) {
	if !v.Connected {
		return "", fmt.Errorf("vehicle %s offline", v.VIN)
	}
	status := "OK"
	if v.FuelPct < 15 {
		status = "LowFuel"
	}
	return fmt.Sprintf("Vehicle %s fuel=%.1f%% loc=%s status=%s", v.VIN, v.FuelPct, v.Location, status), nil
}

func (v Vehicle) Exec() (string, error) {
	// example action: dispatch maintenance if fuel extremely low
	if v.FuelPct < 5 {
		return "", fmt.Errorf("emergency: vehicle %s requires immediate service", v.VIN)
	}
	return fmt.Sprintf("Vehicle %s telemetry sent", v.VIN), nil
}

// ---------------------------
// Polymorphism Demo: process all Reporters and Executors
// ---------------------------

func runReports(reporters []Reporter) {
	fmt.Println("=== RUN REPORTS ===")
	for _, r := range reporters {
		id := r.ID()
		s, err := r.Report()
		if err != nil {
			log.Printf("[REPORT ERROR] %s: %v\n", id, err)
			continue
		}
		fmt.Printf("[REPORT] %s => %s\n", id, s)
	}
}

func runExecutions(executors []Executor) {
	fmt.Println("=== RUN EXECUTIONS ===")
	for _, e := range executors {
		// Exec returns result or domain error
		res, err := e.Exec()
		if err != nil {
			log.Printf("[EXEC ERROR] %v\n", err)
			continue
		}
		fmt.Printf("[EXEC] %s\n", res)
	}
}

// health checker demo using interface type assertion
func runHealthChecks(items []interface{}) {
	fmt.Println("=== RUN HEALTH CHECKS ===")
	for _, it := range items {
		if hc, ok := it.(HealthChecker); ok {
			if err := hc.Check(); err != nil {
				log.Printf("[HEALTH] %T failed check: %v\n", it, err)
			} else {
				fmt.Printf("[HEALTH] %T OK\n", it)
			}
		}
	}
}

func main() {
	// create industry objects
	trades := []Trade{
		{"INFY", 100, 1500, 1600},
		{"TCS", 50, 2000, 1950},
	}

	towers := []Tower{
		{"BLR-001", 75, 1200},
		{"BLR-002", 98, 2300}, // degraded/critical
	}

	saasApps := []SaaSApp{
		{"BillingService", 1200, 60, 45},
		{"Analytics", 300, 92, 88}, // unhealthy
	}

	supplyNodes := []SupplyNode{
		{"WH-DEL-1", 500, 120},
		{"WH-MUM-2", 15, 80}, // low stock
	}

	vehicles := []Vehicle{
		{"WVWZZZ1JZXW000001", 55.0, true, "Munich"},
		{"WVWZZZ1JZXW000002", 4.0, true, "Ingolstadt"}, // emergency fuel
	}

	// Build slices of interfaces (polymorphism)
	var reporters []Reporter
	// add trades, towers, saasApps, supplyNodes, vehicles where they implement Reporter
	for _, tr := range trades {
		reporters = append(reporters, tr)
	}
	for _, tw := range towers {
		reporters = append(reporters, tw)
	}
	for _, s := range saasApps {
		reporters = append(reporters, s)
	}
	for _, sn := range supplyNodes {
		reporters = append(reporters, sn)
	}
	for _, v := range vehicles {
		reporters = append(reporters, v)
	}

	// run reporting stage (polymorphic)
	runReports(reporters)

	// Build executors (some types implement Executor)
	var executors []Executor
	for _, tr := range trades {
		executors = append(executors, tr)
	}
	for _, sn := range supplyNodes {
		executors = append(executors, sn)
	}
	for _, v := range vehicles {
		executors = append(executors, v)
	}

	// run executions
	runExecutions(executors)

	// health checks: some types implement HealthChecker (Tower, SaaSApp)
	var healthItems []interface{}
	for _, tw := range towers {
		healthItems = append(healthItems, tw)
	}
	for _, s := range saasApps {
		healthItems = append(healthItems, s)
	}
	runHealthChecks(healthItems)

	// short pause to simulate streaming batches
	time.Sleep(80 * time.Millisecond)
	fmt.Println("=== ALL INDUSTRY MODULES RUN COMPLETED ===")
}

