package main

import "fmt"

//shopping data
// floor

func main() {

	mall := map[int]map[string]int{
		// floor
		1: {"shops": 5, "restroom": 1, "come-go": 1},
		2: {"shops": 4, "restroom": 2, "kfc": 1},
		3: {"pvr": 5, "food shop": 5, "play area": 1},
	}

	for floor, shops := range mall {
		fmt.Println(floor, shops["shops"], shops["restroom"])
	}

	for i, j := range [5]int{1, 2, 3, 4, 5} {
		fmt.Println(i, j)
	}

	customersmall := map[string]map[string][]string{

		"9900367097": {"pvr": []string{"damka1", "damka2"}},
		"9988899988": {"lifesyle": []string{"tshirt", "pants"}},
	}

	for cphone, shopping := range customersmall {
		fmt.Println(cphone, shopping)
	}
}
