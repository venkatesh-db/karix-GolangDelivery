package main

import "fmt"

// offering to god is different for each devotee
// worhsip --> sweet pongal  , rice , flowers
// worship -->  hen , flowers, rice
// worship -->  milk , flowers, fruits

// prayer --> what we offer is difference is ploymorphism

type Devotee struct {
	devotter string
	culture  string
	values   []string
}

type Dusetti struct {
	Devotee  // inheritence
	offering []string
}

func (p *Dusetti) worship() {

	p.offering = []string{"sweet pongal", "rice", "flowers"}
}

type bele struct {
	Devotee  // inheritence
	offering []string
}

func (b *bele) worship() {
	b.offering = []string{"hen", "flowers", "rice"}
}

type harjani struct {
	Devotee
	decorateflowers []string
	fruits          []string
	sweets          []string
}

func (h *harjani) worship() {

	h.decorateflowers = []string{"milk", "flowers", "fruits"}
	h.sweets = []string{"special sweet"}
	h.fruits = []string{"apple", "banana"}

	fmt.Println(h)
}

func main() {

	var mohak harjani

	// to pray we need to create object
	// allocate memory for struct variable

	mohak.devotter = "shiva"
	mohak.culture = "hindu"
	mohak.values = []string{"kindness", "honesty", "patience"}

	mohak.worship() // call method offering to god

	var sagar bele

	sagar.devotter = "shiva"
	sagar.culture = "hindu"
	sagar.values = []string{"duty", "respect", "love"}

	sagar.worship() // call method offering to god

}
