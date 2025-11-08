package main

import "fmt"

/*
Emotional Logic Summary

if / else → Decision making.

switch → Clean branching.

for → Repetition and rhythm.

*/


func main() {
    var happiness int = 85
    var calmness int = 90

    // Comparison operators
    if happiness > calmness {
        fmt.Println("😊 You are more happy than calm.")
    } else if happiness == calmness {
        fmt.Println("⚖️ Balanced emotions — inner peace.")
    } else {
        fmt.Println("🌿 You are calmer today.")
    }

    // Switch expression
    moodLevel := "excited"

    switch moodLevel {
    case "happy":
        fmt.Println("💖 Keep spreading joy!")
    case "excited":
        fmt.Println("⚡ You’re full of energy today!")
    default:
        fmt.Println("🌸 Stay positive!")
    }

    // Loop for self-improvement
    for day := 1; day <= 3; day++ {
        fmt.Printf("Day %d — Reflect, Learn, Grow\n", day)
    }
}
