package main

import "fmt"

//  fairandloevly
//  fairhandsome

type Wowskin struct {
	cream string
}

type Glowingskin struct {
	brand string
}

// face
type Menface struct {
	Glowingskin // Inheritance
	Layeredskin string
}

type Womenface struct {
	Glowingskin     // Inheritance
	Wowskin         // Inheritance
	Softlayeredskin string
}

func main() {

	var pasupathi Menface

	pasupathi.Layeredskin = "Glowingskin Cream Applied"
	pasupathi.brand = "Layeredskin Applied"

	fmt.Println(pasupathi)

	var seetha Womenface

	seetha.cream = "Glowingskin Cream Applied"
	seetha.Wowskin.cream = "Wowskin Cream Applied"
	seetha.Glowingskin.brand = "Fairandlovely Applied"
	seetha.Softlayeredskin = "Softlayeredskin Applied"

	fmt.Println(seetha)

}
