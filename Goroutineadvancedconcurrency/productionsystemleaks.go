/*

Most real production bugs happen due to incorrect channel usage:

❌ Goroutine waiting forever on <-chan
❌ Channel never closed
❌ Unbuffered channel blocking writers → goroutine stuck
❌ Select with missing timeout
❌ Work queued without backpressure
❌ Fan-out or worker pools without cancellation
❌ Forgetting to drain channels

Goroutine leaks eventually kill production systems
→ memory grows
→ CPU spikes
→ expensive incidents

rule 1 - Always combine unbuffered channels with select + ctx.Done()

rule 2-- Always use buffered channels for work queues

⚠️ Must choose buffer size carefully

Small buffer → writers will block
Large buffer → memory is wasted

✔️rule 3 — Select Statement → The MOST IMPORTANT LEAK PREVENTION TOO



Rule 1 — Every goroutine must have an exit path

Rule 2 — Close channels only by the sender

Rule 3 — Never send to a closed channel

Rule 4 — Drain channels before shutting down

Rule 5 — For queues → always use buffered channels

Rule 6 — Use worker pools for controlled concurrency

*/

package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickets := make(chan int, 10) // 10 buffered slots → 10 tickets

	// fill tickets
	for i := 10; i >= 1; i-- {
		tickets <- i
	}
	close(tickets)

	var wg sync.WaitGroup

	for person := 1; person <= 50; person++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			select {
			case ticket, ok := <-tickets:
				if !ok {
					fmt.Printf("Person %2d → No tickets left\n", id)
					return
				}
				fmt.Printf("Person %2d booked ticket %d\n", id, ticket)

			case <-ctx.Done():
				// prevents leak if system is shutting down
				return
			}

		}(person)
	}

	wg.Wait()
	fmt.Println("All requests processed")
}
