package main

import (
	"fmt"
)

var globalVar int = 42

func heap() {

	var unlimitedmeal []string = []string{"Pizza", "Burger", "Pasta", "Salad", "Sushi"}
	fmt.Println("Heap Memory Segment Example:", unlimitedmeal)

	 selfies :=map[string][]string{
		"mysore":[]string{"place","palace","temple"}, 
		"goa":[]string{"beach","party","fun"},
	 }
	 fmt.Println("Map in Heap Memory Segment Example:", selfies)
}

func main() {
	localVar := "Hello, Go!"
	fmt.Println("Global Variable:", globalVar)
	fmt.Println("Local Variable:", localVar)
	heap()
}
