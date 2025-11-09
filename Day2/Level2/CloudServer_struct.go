package main

import (
	"fmt"
	"log"
)

/*

Real-world cloud monitoring scenario

Clear struct and methods

Logging + alert system

Clean, production-ready

*/

// -------------------------
// Struct: CloudServer
// -------------------------
type CloudServer struct {
	Name        string
	CPUUsage    int // percentage
	MemoryUsage int // percentage
	IsActive    bool
}

// Method: Server Health Status
func (s CloudServer) HealthStatus() string {
	if s.CPUUsage < 75 && s.MemoryUsage < 80 {
		return "Healthy"
	}
	return "Critical"
}

// Method: Alert if server is overloaded
func (s CloudServer) Alert() {
	if s.HealthStatus() == "Critical" {
		log.Printf("[ALERT] Server %s is overloaded! CPU: %d%%, Memory: %d%%\n", s.Name, s.CPUUsage, s.MemoryUsage)
	}
}

func main() {
	servers := []CloudServer{
		{"Server-A", 65, 70, true},
		{"Server-B", 85, 90, true}, // Critical
		{"Server-C", 45, 60, true},
	}

	for _, server := range servers {
		fmt.Printf("[INFO] %s Status: %s\n", server.Name, server.HealthStatus())
		server.Alert()
	}
}
