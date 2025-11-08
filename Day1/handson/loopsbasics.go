

package main

import "fmt"

func main() {

	// vehcile parking --> no of floor pass 2 to 3
	// floor1- 100
	// floor2- 100

	floors := 0 // to park my vehcile
	avacpacityfloor := 0
	invech := 1
	parkav := false
	avacpacityscndfloor := 1

	for floors < 3 { // 3 floors

		if invech <= avacpacityfloor && parkav == true {

			parkedinfirstfloor := true
			fmt.Println("parkedinfirstfloor", parkedinfirstfloor)
			avacpacityfloor = avacpacityfloor - 1
			break // skip to go to second floor

		} else if floors == 1 && avacpacityscndfloor == 1 {

			parkedinsecondfloor := true
			fmt.Println("parkedinsecondfloor", parkedinsecondfloor)
			avacpacityscndfloor = avacpacityscndfloor - 1
			break
		}

		floors = floors + 1
	}

	fmt.Println(floors, avacpacityscndfloor, invech)

}
