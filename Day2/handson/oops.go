package main

import "fmt"

// family -->[dusetti]  mother father sister brother
// mother -> [dusetti]  god --> durgamatha
// sister -> [dusetti]  god ->   hanuman
// brother -> [dusetti] god -->  shiva

type Dusetti struct{

	devotter string
	culture string
	values []string

}

type spiritual struct{
	Dusetti // inheritence
	god string
	temple string
	festival []string
}

func main(){


	var keerthi spiritual  // in order to do prayer we need to create object 
	// when we initalise the object we can tell the ask to god 

	keerthi.culture = "pooja offering"
	keerthi.devotter = "lord shiva"
	keerthi.values = []string{"respecting elders","helping others"}

	keerthi.god = "durgamatha"
	keerthi.temple = "durgamatha temple"
	keerthi.festival = []string{"dusserha","navarathri"}

	var venu spiritual 

	venu.culture = "prayers and rituals"
	venu.devotter = "lord shiva"
	venu.values = []string{"bravery","loyalty"}

	venu.god = "hanuman"
	venu.temple = "hanuman temple"
	venu.festival = []string{"hanuman jayanti"}

	fmt.Println(keerthi)
	fmt.Println(venu)

}
