package main

import "fmt"


// common methods for all chcoloate 
type taste interface {

	choloclatemaking()
	// any chocolate we eat are  unhealthy

}

type Fat struct {
	Cholesterol int
	Sugar       int
	HealthNotes string
}

func (f Fat) paybillstodoctor(ch int, su int, hn string) {

	f.Cholesterol = ch
	f.Sugar = su
	f.HealthNotes = hn
	fmt.Println(f)
}

type Munch struct {

	Fat // inheritence

	brand   string
	flavour []string // polymorphism

}

func (m Munch) choloclatemaking() {

	m.flavour = []string{"Wafer + chocolate coating"}
	for _, flav := range m.flavour {
		println("Munch flavour:", flav)
	}

}

type kitkat struct {

	Fat // inheritence

	variety string
	flavour []string // polymorphism

}

func (k kitkat) choloclatemaking() {

	k.flavour = []string{"Milk chocolate"}

	fmt.Println(k.flavour)
}

func main() {

	// shop owner wants to sell different types of chocolates - interface
	// customer ask munch chocolate -> face venkat -> much object
    // shop ==> collect money ftrom customer
    // shop owner delivery product to customer

	var advertisment taste // remind to eat chocolate
     var venkat kitkat //to  buy the cholclate
	advertisment =&venkat // address of object
	advertisment.choloclatemaking()

	/*
	venkat.paybillstodoctor(2, 47, "High sugar intake may lead to diabetes")
	venkat.choloclatemaking()
	*/
}