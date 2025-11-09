
package main 
import "fmt"

func fearland(inv string,ch chan string){

	// <- read mesg of agent from channel ch
	fmt.Println("fear land of growth issues",<-ch)
}

func main(){
	
	fmt.Println("main started")

     ch:= make(chan string) //heap memory agent to owner and customer

	go fearland("1cr-2cr",ch) // fearland is goroutine owned by venkatesh 
	// venkatesh ownerland --> buyer --> builder , customer 
    // owner --> broker agent --> customer or builder

	//  broker agent --> cost 2 cr 1.9 cr negotiation

	ch<-"2cr" // when customer convey his price immediately owner is notified
	// immediately goroutine fearland is executed

     fmt.Println("you price is notifed to goroutine")

	fmt.Println("main endeed")

}
