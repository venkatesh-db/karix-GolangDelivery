package main

import "fmt"

/*

const = values that never change (emotions remain stable today 😄).

iota = automatically increments values (great for enumerations).

type EmotionScore int = custom type (makes meaning explicit).

*/



const (
    Happiness   = 100
    Calmness    = 80
    Confidence  = 95
)

// Enum-like pattern using iota
const (
    Morning = iota // 0
    Afternoon      // 1
    Evening        // 2
)

type EmotionScore int

func main() {
    var todayMood EmotionScore = Happiness

    fmt.Println("🌞 Time of Day:", Morning)
    fmt.Println("❤️ Mood Score:", todayMood)
    fmt.Println("Confidence Level:", Confidence)
}

