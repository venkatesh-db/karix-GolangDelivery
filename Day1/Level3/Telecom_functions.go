package main

import (
	"errors"
	"fmt"
)

// Function to check signal strength
func checkSignal(strength int) error {
	if strength < -85 {
		return errors.New("weak signal: cannot make call")
	}
	return nil
}

func main() {
	signalLevels := []int{-70, -90, -80} // dBm values

	for i, s := range signalLevels {
		err := checkSignal(s)
		if err != nil {
			fmt.Printf("Tower %d ❌ %s\n", i+1, err)
		} else {
			fmt.Printf("Tower %d ✅ Signal is good\n", i+1)
		}
	}
}

