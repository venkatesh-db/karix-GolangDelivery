package main

import "fmt"

type templepcore struct {
	nameofgod  string
	face       string
	architect  string
	templeType string
	sitted     string
}

func (t *templepcore) freedarshan() {

	t.nameofgod = "vishnu"
	t.face = "four faced"
	t.architect = "dravidian style"
	t.templeType = "temple on land"
	t.sitted = "on lotus"
	fmt.Println("famoius temple", t)

}

func main() {

	// To see god

	// create object

	var nayan templepcore // stack

	var rokesh *templepcore = new(templepcore) // heap.

	// assign values to the fields

	nayan.nameofgod = "vishnu"
	nayan.face = "four faced"
	nayan.architect = "dravidian style"
	nayan.templeType = "temple on land"
	nayan.sitted = "on lotus"

	fmt.Println("direct view local temple ", nayan)

	nayan.freedarshan()  // & nayan
	rokesh.freedarshan() // rokesh

}
