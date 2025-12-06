package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {

		for i := 0; i < 5; i++ {

			select {

			case <-ctx.Done():
				fmt.Println("worker 1 cancelled")
				return ctx.Err()
			default:
				// Do some work here
				fmt.Println("first goroutine")
				fmt.Println("worker default")
				time.Sleep(300 * time.Millisecond)
			}

		}
		return nil

	})

	g.Go(func() error {

		fmt.Println("second gororutine")
		time.Sleep(1 * time.Second)
		fmt.Println("second gororutine completed")
		return errors.New("worker 2 failed")

	})

	fmt.Println("first wait for gororutine")

	if err := g.Wait(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("All workers completed successfully")

	}

}
