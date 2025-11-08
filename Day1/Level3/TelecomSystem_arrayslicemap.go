package main

import "fmt"

func main() {
	// Fixed towers
	towers := [3]string{"BLR-KA-001", "BLR-KA-002", "BLR-KA-003"}
	fmt.Println("[TELECOM] Towers Array:", towers)

	// Active connections slice
	activeConnections := []int{120, 230, 95}
	activeConnections = append(activeConnections, 150)
	fmt.Println("[TELECOM] Active Connections Slice:", activeConnections)

	// Map — TowerID → SignalStrength
	signalStrength := map[string]float32{
		"BLR-KA-001": -75.5,
		"BLR-KA-002": -80.0,
	}
	signalStrength["BLR-KA-003"] = -70.2
	fmt.Println("[TELECOM] Signal Map:", signalStrength)
}
