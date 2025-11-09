package main

import "fmt"

// Define a simple struct for a person
type Person struct {
	Name    string
	Age     int
	Country string
}

// Method assigned to struct
func (p Person) Greet() {
	fmt.Printf("Hello, my name is %s. I am %d years old from %s.\n", p.Name, p.Age, p.Country)
}

func main() {
	// Initializing struct
	person1 := Person{"Venkatesh", 39, "India"}
	person1.Greet()
}
