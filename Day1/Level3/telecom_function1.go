
package main

import (
	"errors"
	"fmt"
)

// Function: monitor tower and generate alert
func monitorTower(tower string, signal int) (string, error) {
	if signal < -85 {
		return "", errors.New("signal below threshold: alert raised")
	}
	return fmt.Sprintf("%s signal is healthy (%ddBm)", tower, signal), nil
}

func main() {
	towers := map[string]int{
		"TowerA": -70,
		"TowerB": -90,
		"TowerC": -80,
	}

	for tower, signal := range towers {
		status, err := monitorTower(tower, signal)
		if err != nil {
			fmt.Printf("❌ %s: %s\n", tower, err)
		} else {
			fmt.Printf("✅ %s\n", status)
		}
	}
}
