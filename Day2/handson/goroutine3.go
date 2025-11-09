package main 

import "fmt"


func babycrying(name string,stopcry chan int){

	fmt.Println("baby is crying....",name)

	for{ // non stop crying until stopped by mother
		item:=<-stopcry
		fmt.Println("baby stopped crying on getting ",item)
	}

	fmt.Println("baby is end....")


}

func main(){

	stopcry:= make(chan int) //heap memory agent to father and baby

	go babycrying("mother",stopcry) // babycrying is goroutine owned by mother


	for i:=0;i<5;i++{
		fmt.Println("gto stop crying" ,i)
        stopcry <- i // when mother convey item immediately baby is notified
		// immediately goroutine babycrying is executed
	}

	fmt.Println("main endeed")
}