package main

import (
	"fmt"
	"time"
)

// brain --> concurrency model
// think1 --> regarding work
// think2 --> personal life land issues
// think3 --> buy rich house big car big salary
// think4 --> love mother wife children

// Regarding work --> function's actvity1-start end
// Personal life land issues --> running 5 year's problem goroutine
// Think3 --> buy rich house big car big salary --> function's
//  daily running --> buy rich house big car big salary --> goroutine
//  Think4 --> started love baby when  born is running --> not ended ,started started 5 years ago

func lovebaby(born string) {

	fmt.Println("baby born started")

	fmt.Println("baby end ")

} // end of the day loving baby is not stopped

func landissues(year string) {

	fmt.Println("land issues started")

	fmt.Println("land issues startendeded")

} // end of the day or end of month issue resolved

func main() {

	fmt.Println("main started")

	//lovebaby("land issues") // function call executes function
	//landissues("5 years")

	go lovebaby("land issues") // lovebaby is goroutine go lovebaby("land issues" )
	go landissues("5 years")   // landissues is goroutine

	// inorder to execute a gororutine we use time.Sleep method
	time.Sleep(2 * time.Second) // main goroutine is sleeping for 2 seconds to wait other goroutines to complete

	fmt.Println("main ended")
}
