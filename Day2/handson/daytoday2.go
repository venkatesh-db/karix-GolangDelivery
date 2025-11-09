package main

// day1- memory managment
// day2 - reusable oops 

// complement human or recognsing team or people
// comment    human
// fun        jokes 
// deadline of work 
// stress in job 
// relax in office    cofee tea 
// games       

import "fmt"

// 1. desiging 
// 2. memory single continous slice map 
// 3. srp principal 


// complement human or recognsing team or people  --> interface 
// comment    human -->   function human1 -function human2
// fun        jokes  -->  struct 
// deadline of work  -->  interface
// stress in job     -->  interface , map 
// relax in office    cofee tea --> interface
// games             -->  interface

func suresh(fun string){

  scolding:= ramesh(fun)

  	fmt.Println(scolding)
}

func ramesh(fun string) string{

	fmt.Println(fun)

	return "jamun boy"

}

func main(){

	suresh("kesari boy")

	var commonfacility board
	var hemath board1
	commonfacility = &hemath
	commonfacility.width()

}

type  board interface {

	 width()
}

type board1 struct {
	 length int
}
	
func (b board1) width(){
	fmt.Println("length of board is ",b.length)
}

type board2 struct {
	  breadth int
}

func (b board2) width(){
	fmt.Println("breadth of board is ",b.breadth)
}	
