package main

import (
	"fmt"
	"sync"
	"time"
)

var tickets = 10
var mu sync.Mutex

func bookTicket(person int, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(person%5) * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if tickets > 0 {
		fmt.Printf("Person %2d booked ticket number %d\n", person, tickets)
		tickets--
	} else {
		fmt.Printf("Person %2d found no tickets\n", person)
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go bookTicket(i, &wg)
	}

	wg.Wait()
	fmt.Println("Final tickets:", tickets)
}
