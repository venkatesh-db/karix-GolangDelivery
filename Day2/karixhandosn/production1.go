package main 

import "fmt"

type Classmate struct {
	Name  string
}

type  Roommate struct {

	 Classmate    // One object integral part of another struct
    // base - 100 bytes + derived - 200 bytes = 300 bytes


	 ram Classmate // One object integral part of another object
	   // base - 100 bytes + derived - 200 bytes = 300 bytes

	 agge *Classmate // pointer to another struct
	   // base  ponly one pointer to 100 bytes 

    rent int

}


func main() {

	fmt.Println("Production code ")


	var venkatesh Roommate

	venkatesh.Name = "venkatesh"
    venkatesh.rent = 5000

	venkatesh.ram.Name = "sairam" // composition 

	venkatesh.agge = &Classmate{Name: "agg"} // aggregation


}



