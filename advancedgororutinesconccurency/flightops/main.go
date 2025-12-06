package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/oklog/run"
	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logger := log.New(os.Stdout, "[flightops] ", log.LstdFlags|log.Lmicroseconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Println("context-managed flight monitors")
	flightWatchWithContext(ctx, logger)

	logger.Println("crew comms via channels")
	crewChannelDemo(logger)

	logger.Println("seat inventory guarded by mutex")
	seatInventoryDemo(logger)

	logger.Println("maintenance errors propagated over channel")
	maintenanceFanout(ctx, logger)

	logger.Println("flight readiness checks with errgroup")
	flightPreparationErrgroup(ctx, logger)

	logger.Println("airline lifecycle orchestrated by oklog/run")
	airlineLifecycle(logger)
}

func flightWatchWithContext(parent context.Context, logger *log.Logger) {
	ctx, cancel := context.WithTimeout(parent, 2500*time.Millisecond)
	defer cancel()

	flights := []string{"AI101", "AI447", "AI907"}
	var wg sync.WaitGroup

	for i, code := range flights {
		wg.Add(1)
		go func(idx int, flight string) {
			defer wg.Done()
			if err := monitorFlight(ctx, flight, time.Duration(400+idx*400)*time.Millisecond); err != nil {
				logger.Printf("monitor %s halted: %v", flight, err)
				return
			}
			logger.Printf("monitor %s completed", flight)
		}(i, code)
	}

	wg.Wait()
}

func monitorFlight(ctx context.Context, flight string, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s canceled: %w", flight, ctx.Err())
	}
}

func crewChannelDemo(logger *log.Logger) {
	briefing := make(chan string)
	go func() { briefing <- "Pushback approved" }()
	logger.Printf("tower -> captain: %s", <-briefing)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	announcements := make(chan string, 3)
	go func() {
		defer close(announcements)
		msgs := []string{"Cabin secure", "Cross-check complete", "Doors armed"}
		for _, msg := range msgs {
			select {
			case <-ctx.Done():
				return
			case announcements <- msg:
				logger.Printf("cabin crew queued: %s", msg)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("crew channel timeout: %v", ctx.Err())
			return
		case msg, ok := <-announcements:
			if !ok {
				logger.Printf("crew announcements drained")
				return
			}
			logger.Printf("flight deck received: %s", msg)
		}
	}
}

func seatInventoryDemo(logger *log.Logger) {
	type seatInventory struct {
		sync.Mutex
		available int
	}

	s := &seatInventory{available: 180}
	var wg sync.WaitGroup

	reserve := func(count int) {
		s.Lock()
		defer s.Unlock()
		if s.available >= count {
			s.available -= count
		}
	}

	for _, pax := range []int{2, 4, 1, 3, 5} {
		wg.Add(1)
		go func(seats int) {
			defer wg.Done()
			for i := 0; i < seats; i++ {
				reserve(1)
				time.Sleep(10 * time.Millisecond)
			}
		}(pax)
	}

	wg.Wait()
	logger.Printf("seats remaining: %d", s.available)
}

func maintenanceFanout(ctx context.Context, logger *log.Logger) {
	errCh := make(chan error, 3)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for _, system := range []string{"hydraulics", "avionics", "fuel"} {
			go func(component string) {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
				case <-time.After(time.Duration(200+rand.Intn(200)) * time.Millisecond):
					if component == "avionics" {
						errCh <- fmt.Errorf("%s self-test failed", component)
						return
					}
					errCh <- nil
				}
			}(system)
		}
	}()

	var firstErr error
	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil && firstErr == nil {
			firstErr = err
		}
	}
	<-done

	if firstErr != nil {
		logger.Printf("maintenance halted: %v", firstErr)
		return
	}
	logger.Printf("maintenance cleared all systems")
}

func flightPreparationErrgroup(ctx context.Context, logger *log.Logger) {
	checks := []string{"weather", "payload", "baggage", "fuel"}
	g, ctx := errgroup.WithContext(ctx)

	for _, check := range checks {
		name := check
		g.Go(func() error {
			time.Sleep(time.Duration(150+rand.Intn(200)) * time.Millisecond)
			if name == "baggage" {
				return fmt.Errorf("%s belt jammed", name)
			}
			logger.Printf("%s check passed", name)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Printf("flight readiness aborted: %v", err)
		return
	}
	logger.Printf("flight readiness complete")
}

func airlineLifecycle(logger *log.Logger) {
	var g run.Group
	opsCtx, stopOps := context.WithCancel(context.Background())
	crewCtx, stopCrew := context.WithCancel(context.Background())

	g.Add(func() error {
		logger.Printf("ops control online")
		<-opsCtx.Done()
		logger.Printf("ops control shutting down")
		return opsCtx.Err()
	}, func(err error) {
		logger.Printf("ops interrupt: %v", err)
		stopOps()
	})

	g.Add(func() error {
		logger.Printf("crew scheduler dispatching")
		<-crewCtx.Done()
		logger.Printf("crew scheduler stood down")
		return crewCtx.Err()
	}, func(err error) {
		logger.Printf("crew scheduler interrupt: %v", err)
		stopCrew()
	})

	g.Add(func() error {
		time.Sleep(1 * time.Second)
		return errors.New("airport blackout")
	}, func(error) {
		logger.Printf("infrastructure team engaged")
	})

	if err := g.Run(); err != nil {
		logger.Printf("airline lifecycle exit: %v", err)
	}
}
