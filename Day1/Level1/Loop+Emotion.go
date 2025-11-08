package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	totalRounds := 7

	fmt.Println("🙏 Visiting Temple – Performing 7 Rounds for Blessings\n")

	for round := 1; round <= totalRounds; round++ {
		fmt.Printf("🕉️ Round %d: ", round)

		switch round {
		case 1, 2:
			fmt.Println("Hope for Job – ✅ Expectation met!")
		case 3, 4:
			fmt.Println("Hope for House – ✅ Expectation met!")
		case 5:
			fmt.Println("Possibility – 🚗 BMW Car Offer received (10 Gold equivalent) 💎")
		case 6:
			fmt.Println("Possibility – 🏡 Unexpected blessing: New plot offer 🌿")
		case 7:
			fmt.Println("🙏 Unexpected Outcome – Divine Timing, not yet granted but faith continues 🌸")
		default:
			fmt.Println("🌼 Peaceful round with gratitude")
		}
	}
}



func parks() {
	// park visits in a week
	totalVisits := 2 // visiting park 2 times a week

	for visit := 1; visit <= totalVisits; visit++ {
		fmt.Printf("\n🏞️ Visit %d to the park:\n", visit)

		// each visit has 2 rounds
		for round := 1; round <= 2; round++ {
			if round == 1 {
				fmt.Printf("  Round %d - Ravi Bala walking: 😅 Unexpected Outcome\n", round)
			} else if round == 2 {
				fmt.Printf("  Round %d - Ravi Kiran running: ✅ Expected Outcome\n", round)
			} else {
				fmt.Printf("  Round %d - Cooling down 🧘‍♂️\n", round)
			}
		}
	}
}





func winn() {
	rand.Seed(time.Now().UnixNano())
	outcomes := []string{
		"✅ Job secured",
		"🏡 House registration done",
		"🚗 BMW car offer received",
		"💎 Gold value increased",
		"🌸 Unexpected Outcome – patience needed",
	}

	for round := 1; round <= 7; round++ {
		result := outcomes[rand.Intn(len(outcomes))]
		fmt.Printf("Round %d: %s\n", round, result)
	}
}
