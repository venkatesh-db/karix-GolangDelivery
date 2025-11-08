package main

import "fmt"

func main() {
    // == Mindset: Dress well every day
    mindset := "dress_well"
    if mindset == "dress_well" {
        fmt.Println("🧠 Positive mindset — you look confident today!")
    }

    // != Mom says: Don’t go outside today
    goOutside := false
    if goOutside != true {
        fmt.Println("👩‍👧 Mom: Stay home today, it’s raining outside.")
    }

    // < Marriage gold comparison
    motherGold := 10  // grams
    brideGold := 100  // grams
    if motherGold < brideGold {
        fmt.Println("💍 Bride has more gold for the wedding!")
    }

    // > Father proud of IIT son
    fatherEducation := "BTech"
    sonEducation := "IIT"
    if sonEducation > fatherEducation { // symbolic emotional comparison
        fmt.Println("👨‍👦 Father: My son studied at IIT, I’m proud of him!")
    }

    // > Salary comparison between daughter and someone
    daughterSalary := 30_00_000  // 30 lakhs
    neighborDaughterSalary := 20_00_000
    if daughterSalary > neighborDaughterSalary {
        fmt.Println("💼 My daughter earns more — 30L vs 20L!")
    }

    // > Son at USA vs working in India
    mySonLocation := "USA"
    yourSonLocation := "India"
    if mySonLocation > yourSonLocation { // symbolic compare (alphabetically)
        fmt.Println("🌎 My son works in USA — proud parent moment!")
    }
}

