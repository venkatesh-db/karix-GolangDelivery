package main

import "fmt"

// 1 line of code
var smile int8 = 25 // Data decide memory - 1 byte i8

// smile --> variable scope exist in this file
//  lieftime --> memory exist till end of the program

func main() {

	var nofcars int8 = 2

	// nofcars --> variable scope exist in this block
	// lieftime --> memory exist till end of the block

	smiles := 25 // Decide memory -int8

	fmt.Println(smiles)
	fmt.Println(nofcars)

} // Scope & lifetime
