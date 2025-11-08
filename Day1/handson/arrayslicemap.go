
// supermartket
// medical shop

package main

import "fmt"

func medicalshop() {

	// tabletname scrpies - fixed static memory

	var doctorpresciption [2]string = [2]string{"dolo", "syrcpg"}

	var medicinebill [2][5]string = [2][5]string{[5]string{"cipla", "local1", "local2", "lcoal3", "local4"},
		[5]string{"dolo", "solo", "volo", "zolio", "sunlo"},
	}

	// dynamic memory add n no of medicie

	bills := map[string]string{
		"cipla":  "dolo",
		"local1": "solo",
	}
	fmt.Println(bills)

	bills["xmen"] = "sugarless tablest"
	bills["throuht infection"] = "stepslis"
	fmt.Println(bills)
	fmt.Println(doctorpresciption, medicinebill)

}

func supermarket() {

	// category1 -continous 5 memory
	var foodgrocies [5]string = [5]string{"dal", "chawal", "milk", "curd", "sugar"}
	// category2 -continous 5 memory
	var houstems [5]string = [5]string{"stck", "soap", "chair", "fans", "bulbs"}

	fmt.Println(foodgrocies, houstems)

	// 2 categories -->continous 5 memory [5]string
	var totalitems [2][5]string = [2][5]string{[5]string{"dal", "chawal", "milk", "curd", "sugar"},
		[5]string{"stck", "soap", "chair", "fans", "bulbs"}}

	fmt.Println(totalitems)
}

func main() {

	supermarket()
	medicalshop()
}
