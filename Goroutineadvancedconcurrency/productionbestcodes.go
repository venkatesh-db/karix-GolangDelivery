
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// Cancel everything if no one wins within 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	job := make(chan string, 1) // Buffered → prevents blocking
	results := make(chan string, 10)

	var wg sync.WaitGroup

	for person := 1; person <= 10; person++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// Simulate person trying to grab the job
			select {
			case job <- fmt.Sprintf("Candidate %d got the job!", id):
				// This person succeeded
				results <- fmt.Sprintf("Candidate %d → SELECTED", id)
				cancel() // cancel all other candidates

			case <-ctx.Done():
				// Could not get the job, or job already taken
				results <- fmt.Sprintf("Candidate %d → NOT selected", id)
				return
			}

		}(person)
	}

	// Close results after all candidates finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Display final output
	for r := range results {
		fmt.Println(r)
	}

	fmt.Println("Hiring cycle complete.")
}
