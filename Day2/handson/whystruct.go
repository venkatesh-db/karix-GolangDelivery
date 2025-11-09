package main

import "fmt"

type  memory struct{

	// single memory
	parkingcars int // not decided memory

	// continous  fixed memory 
	officebuildings [5]string //  not decide memory

	// continous dynamic memory- slice
	projectteammembers []string
	// heap memory
	
	// continous infinite memory - map //  not decide memory
	tecpark map[string][]string

}


func main() {

// single memory 
var parkingcars int= 500 // stack memory

// continous  fixed memory 
var officebuildings [5]string= [5]string{"A","B","C","D","E"} // stack memory

// continous dynamic memory- slice

var projectteammembers []string= []string{"suman","ajith","harish","seetha"}
// heap memory
projectteammembers= append(projectteammembers,"vijay") 

// continous infinite memory - map // heap memory
var tecpark map[string][]string= map[string][]string{ 
	"blockA":[]string{"google","microsoft","amazon"},
	"blockB":[]string{"tcs","infosys","wipro"},
}
tecpark["blockC"]= []string{"cognizant","hcl","ibm"}


var venkat memory // stack memory allocate memory
// initalise struct members
venkat.parkingcars= 300
venkat.officebuildings= [5]string{"F","G","H","I","J"}
venkat.projectteammembers= []string{"ram","lakshman","bharath"}
venkat.projectteammembers= append(venkat.projectteammembers,"sita")
venkat.tecpark= map[string][]string{
	"blockX":[]string{"meta","netflix","spotify"},
}
venkat.tecpark["blockY"]= []string{"adobe","salesforce"}

fmt.Println("parkingcars:",parkingcars)
fmt.Println("officebuildings:",officebuildings)
fmt.Println("projectteammembers:",projectteammembers)
fmt.Println("tecpark:",tecpark)

fmt.Println("venkat parkingcars:",venkat.parkingcars)
fmt.Println("venkat officebuildings:",venkat.officebuildings)
fmt.Println("venkat projectteammembers:",venkat.projectteammembers)
fmt.Println("venkat tecpark:",venkat.tecpark)
}