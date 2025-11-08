
package main

/*
Production Insight:

Function returns error instead of just printing

Reusable → multiple bookings

Conditional error handling → standard practice in real Go services

*/


import (
	"errors"
	"fmt"
)

// Function: bookTrain returns error if seats unavailable
func bookTrain(pnr string, seatAvailable bool) error {
	if seatAvailable {
		return nil // booking successful
	}
	return errors.New("seat not available")
}

func main() {
	bookings := map[string]bool{
		"PNR123": true,
		"PNR124": false, // seat not available
		"PNR125": true,
	}

	for pnr, available := range bookings {
		err := bookTrain(pnr, available)
		if err != nil {
			fmt.Printf("[ERROR] Booking failed for %s: %s\n", pnr, err)
		} else {
			fmt.Printf("[SUCCESS] Booking successful for %s ✅\n", pnr)
		}
	}
}
