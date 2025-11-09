
package main

import (
	"fmt"
	"log"
	"time"
)

/*

We make a Monitor interface and two implementations (CloudMonitor, TelecomMonitor).

Use pointer receivers for methods that might mutate state; value receivers for read-only.

Show interface slice processing (polymorphism) and basic error-return pattern.

*/


// Monitor is a production-style interface used by monitoring pipelines.
type Monitor interface {
	// ID returns a stable identifier for logs and metrics
	ID() string
	// Snapshot returns a human-readable snapshot or error
	Snapshot() (string, error)
	// Alert generates an alert if needed (non-blocking)
	Alert()
}

// CloudMonitor implements Monitor for cloud servers.
type CloudMonitor struct {
	Name      string
	CPU       int // %
	Memory    int // %
	IsPrimary bool
}

// ID uses value receiver (no mutation)
func (c CloudMonitor) ID() string { return "cloud:" + c.Name }

// Snapshot reads fields and may return an error if unhealthy
func (c CloudMonitor) Snapshot() (string, error) {
	if c.CPU > 95 || c.Memory > 95 {
		return "", fmt.Errorf("critical: %s overloaded (cpu=%d,memory=%d)", c.Name, c.CPU, c.Memory)
	}
	return fmt.Sprintf("Cloud %s OK (cpu=%d,memory=%d)", c.Name, c.CPU, c.Memory), nil
}

// Alert logs alerts (pointer receiver not required)
func (c CloudMonitor) Alert() {
	if c.CPU > 85 || c.Memory > 85 {
		log.Printf("[ALERT] Cloud %s crossing threshold cpu=%d memory=%d\n", c.Name, c.CPU, c.Memory)
	}
}

// TelecomMonitor implements Monitor for telecom towers.
type TelecomMonitor struct {
	TowerID string
	Load    int // %
	Users   int
}

// ID implements Monitor
func (t TelecomMonitor) ID() string { return "tower:" + t.TowerID }

// Snapshot implements Monitor
func (t TelecomMonitor) Snapshot() (string, error) {
	if t.Load >= 100 || t.Users < 0 {
		return "", fmt.Errorf("invalid metrics for %s (load=%d users=%d)", t.TowerID, t.Load, t.Users)
	}
	return fmt.Sprintf("Tower %s Load=%d Users=%d", t.TowerID, t.Load, t.Users), nil
}

// Alert implements Monitor
func (t TelecomMonitor) Alert() {
	if t.Load > 90 {
		log.Printf("[ALERT] Tower %s overloaded load=%d\n", t.TowerID, t.Load)
	}
}

func main() {
	// Compose monitors polymorphically
	monitors := []Monitor{
		CloudMonitor{"srv-us-east-1", 60, 40, true},
		CloudMonitor{"srv-eu-1", 88, 86, false},
		TelecomMonitor{"BLR-KA-001", 75, 1200},
		TelecomMonitor{"BLR-KA-002", 92, 1500},
	}

	for _, m := range monitors {
		id := m.ID()
		snap, err := m.Snapshot()
		if err != nil {
			log.Printf("[ERROR] %s snapshot error: %v\n", id, err)
			// still call Alert to trigger downstream flows
			m.Alert()
			continue
		}
		// production style: structured log (simple fmt here)
		fmt.Printf("[INFO] %s => %s\n", id, snap)
		m.Alert()
		// small pause to mimic streaming processing
		time.Sleep(50 * time.Millisecond)
	}
}
