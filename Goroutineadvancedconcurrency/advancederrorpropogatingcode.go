package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/oklog/run"
)

func main() {
	var g run.Group

	// Root context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	//
	// --------------------------------------------------------
	// 1. SERVICE GOROUTINE (Example: worker / processor)
	// --------------------------------------------------------
	//
	{
		ctxWorker, stopWorker := context.WithCancel(ctx)

		g.Add(func() error {
			fmt.Println("Worker started...")

			// Simulate doing work until context is canceled
			for {
				select {
				case <-ctxWorker.Done():
					fmt.Println("Worker shutting down...")
					return nil
				default:
					fmt.Println("Worker processing task...")
					time.Sleep(500 * time.Millisecond)
				}
			}
		}, func(error) {
			// Cleanup
			fmt.Println("Stopping worker...")
			stopWorker()
		})
	}

	//
	// --------------------------------------------------------
	// 2. ERROR GENERATOR GOROUTINE (simulate internal failure)
	// --------------------------------------------------------
	//
	{
		g.Add(func() error {
			time.Sleep(3 * time.Second)
			fmt.Println("❌ Internal error occurred in service")
			return errors.New("simulated internal failure")
			// this code  will invoke    Stopping worker and Stopping error generator.

		}, func(error) {

			fmt.Println("Stopping error generator...")
			// Cancel root context
			cancel()
		})
	}

	//
	// --------------------------------------------------------
	// 3. OS SIGNAL HANDLER (CTRL+C)
	// --------------------------------------------------------
	//
	{
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		g.Add(func() error {
			select {
			case <-signals:
				fmt.Println("⚠ Interrupt received")
				return errors.New("interrupt received")
			case <-ctx.Done():
				return nil
			}
		}, func(error) {
			cancel()
		})
	}

	//
	// --------------------------------------------------------
	// START EVERYTHING
	// --------------------------------------------------------
	//

	fmt.Println("Starting service group...")
	err := g.Run()
	fmt.Println("Group shut down:", err)
}
