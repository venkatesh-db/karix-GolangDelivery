
package main

import "fmt"

func main() {
    // 1️⃣ Ground to Earth: < Other Humans
    humility := 80
    otherHumans := 100
    if humility < otherHumans {
        fmt.Println("🌍 Stay grounded — respect others more than yourself.")
    }

    // 2️⃣ Gratitude to God: == All I Have
    gratitude := "everything_i_have"
    allIHave := "everything_i_have"
    if gratitude == allIHave {
        fmt.Println("🙏 Gratitude to God — content with what I have.")
    }

    // 3️⃣ Ego or Attitude: >=
    egoLevel := 95
    humilityLevel := 60
    if egoLevel >= humilityLevel {
        fmt.Println("⚠️ Ego rising — time to balance with humility.")
    }

    // 4️⃣ Respecting: >
    respectGiven := 90
    respectReceived := 70
    if respectGiven > respectReceived {
        fmt.Println("💖 You gave more respect than you received — that’s true character.")
    }

    // 5️⃣ Dissatisfaction
    satisfaction := 40
    desiredHappiness := 100
    if satisfaction < desiredHappiness {
        fmt.Println("😔 Feeling dissatisfied — but room to grow and improve.")
    }
}

