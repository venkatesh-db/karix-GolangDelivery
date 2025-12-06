package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	totalTickets = 10
	totalPeople  = 50
	workerCount  = 10
)

type BookingResult struct {
	PersonID int
	Ticket   int
	WorkerID int
	Status   string // "booked" or "no_ticket" or "canceled"
	Err      error
}

func worker(
	ctx context.Context,
	workerID int,
	persons <-chan int,
	tickets <-chan int,
	results chan<- BookingResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Context canceled (shutdown / timeout)
			return

		case personID, ok := <-persons:
			if !ok {
				// No more people in queue
				return
			}

			// Try to get a ticket for this person
			select {
			case <-ctx.Done():
				// Context canceled while waiting for ticket
				return

			case ticket, ok := <-tickets:
				if !ok {
					// Tickets sold out
					results <- BookingResult{
						PersonID: personID,
						WorkerID: workerID,
						Status:   "no_ticket",
					}
					continue
				}

				// Successfully booked
				results <- BookingResult{
					PersonID: personID,
					Ticket:   ticket,
					WorkerID: workerID,
					Status:   "booked",
				}
			}
		}
	}
}

func main() {
	// Global timeout for the whole operation (defensive)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1. Create tickets channel and pre-fill 10 tickets
	tickets := make(chan int, totalTickets)
	for i := totalTickets; i >= 1; i-- {
		tickets <- i
	}
	close(tickets) // No more tickets will be added

	// 2. Create persons channel (input queue) and results channel
	persons := make(chan int, totalPeople)
	results := make(chan BookingResult, totalPeople)

	var wg sync.WaitGroup

	// 3. Start a fixed-size worker pool
	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(ctx, w, persons, tickets, results, &wg)
	}

	// 4. Enqueue all people into the queue
	go func() {
		defer close(persons)
		for p := 1; p <= totalPeople; p++ {
			select {
			case <-ctx.Done():
				return
			case persons <- p:
			}
		}
	}()

	// 5. Close results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Collect and print results
	bookedCount := 0
	noTicketCount := 0

	for res := range results {
		switch res.Status {
		case "booked":
			bookedCount++
			fmt.Printf("Person %2d booked ticket %2d (handled by worker %2d)\n",
				res.PersonID, res.Ticket, res.WorkerID)

		case "no_ticket":
			noTicketCount++
			fmt.Printf("Person %2d → no tickets left (worker %2d)\n",
				res.PersonID, res.WorkerID)

		case "canceled":
			fmt.Printf("Person %2d → canceled due to shutdown (worker %2d)\n",
				res.PersonID, res.WorkerID)
		}
	}

	fmt.Println("-------------------------------------------------")
	fmt.Printf("Total booked tickets : %d\n", bookedCount)
	fmt.Printf("People without ticket: %d\n", noTicketCount)
	fmt.Printf("Expected booked      : %d\n", totalTickets)

	if err := ctx.Err(); err != nil && err == context.DeadlineExceeded {
		fmt.Println("⚠ Context deadline exceeded (system was too slow or stuck)")
	} else {
		fmt.Println("✅ All requests processed cleanly")
	}
}
