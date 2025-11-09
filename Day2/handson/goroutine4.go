package main

import (
	"fmt"
	"time"
)

func lover(name string, lovech chan string, breakch chan string) {
	
	for i := 1; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)

		if i == 4 {
			breakch <- fmt.Sprintf("%s: sweet heart giving u halwa", name)
			close(breakch)
			return
		}

		lovech <- fmt.Sprintf("%s: lovely deeply waiting new gift", name)
	}
	close(lovech)
}

func moveon(done chan bool) {

	for i := 1; i <= 2; i++ { // recover in 2 months 
		fmt.Println("super man recovering", i)
		time.Sleep(1 * time.Second)
	}
	done <- true
	fmt.Println("move on end")
}

func main() {
	fmt.Println("main entry")

	lovech := make(chan string)
	breakch := make(chan string)
	done := make(chan bool)

	go lover("reetha", lovech, breakch)

	for {
		select {
		case msg, ok := <-lovech: // wait for lover rely 
			if !ok {
				lovech = nil
				continue
			}
			fmt.Println(msg)

		case msg := <-breakch:
			fmt.Println(msg)
			go moveon(done)
			breakch = nil

		case <-done:
			fmt.Println("move on")
			return

		case <-time.After(5 * time.Second):
			fmt.Println("time heals, new girl 😄")
			return
		}
	}
}
