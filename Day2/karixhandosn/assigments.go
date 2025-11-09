// MEemory managment
// Memory oeprations in go lang
// OOPS concepts in go lang

//1.bookmyshow

package main

import "fmt"

func principal() {

	var pointer *string = new(string)
	*pointer = "sairam"

	fmt.Println("value at pointer", *pointer)
	fmt.Println("address at pointer", pointer)

}

type seates struct {
	layout map[string]map[string]map[string]int
	// theatre  // 9:30    // A1 360
}

func bookmyshow() {

	movieData := make(map[string]map[string]seates)
	// 4dx      mon
	movieData["4dx"] = make(map[string]seates)

	movieData["4dx"]["mon"] = seates{layout: make(map[string]map[string]map[string]int)}

	movieData["4dx"]["mon"].layout["9:30"] = make(map[string]map[string]int)

	movieData["4dx"]["mon"].layout["9:30"]["A1"] = make(map[string]int)

	movieData["4dx"]["mon"].layout["9:30"]["A1"]["prime"] = 1

	fmt.Println(movieData)

}

func main() {

	bookmyshow()

	principal()
}

//2. where is my train

// 3. amazon

// 4. whatsupp

// 5. instagram

// 6. swiggy || redbus

// 7. truecaller
