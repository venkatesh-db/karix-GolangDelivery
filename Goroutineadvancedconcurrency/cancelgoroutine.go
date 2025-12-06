package main

import (
	"context"
	"fmt"
	"time"
)

func worker( ctx context.Context) {

	fmt.Println("worker started")

   for {
	   select {
	   case <- ctx.Done():
		   fmt.Println("worker received cancel signal")
		   return
		
	   }
   }

	fmt.Println("worker ended")

}


func main(){
  
 ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx )

	time.Sleep(2 * time.Second)

	cancel() //send cancellation signal to worker

	time.Sleep(1 * time.Second) // give goroutine time to print message

	fmt.Println("main  ended")
}
