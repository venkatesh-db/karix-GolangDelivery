package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oklog/run"
	"golang.org/x/sync/errgroup"
)

const (
	autoLogoutWindow      = 2 * time.Minute
	demoAutoLogout        = 9 * time.Second
	totalConcurrentLogins = 1_000_000
	workerCount           = 256
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logger := log.New(os.Stdout, "[irctc] ", log.LstdFlags|log.Lmicroseconds)

	logger.Println("session auto-logout with context deadlines")
	sessionAutoLogoutDemo(logger)

	logger.Println("seat quota guarded with mutex")
	seatQuotaDemo(logger)

	logger.Println("login flood handled via worker pool + channels")
	loginFloodSimulation(logger)

	logger.Println("ticket pipeline with errgroup fail-fast semantics")
	ticketPipeline(logger)

	logger.Println("control-plane lifecycle using oklog/run")
	controlPlane(logger)
}

func sessionAutoLogoutDemo(logger *log.Logger) {
	sessionCtx, cancel := context.WithTimeout(context.Background(), demoAutoLogout)
	defer cancel()

	heartbeat := make(chan string, 1)
	go simulateUserActivity(sessionCtx, heartbeat, logger)

	logger.Printf("session window: %v (production %v)", demoAutoLogout, autoLogoutWindow)

	for {
		select {
		case <-sessionCtx.Done():
			logger.Printf("session auto-logout triggered: %v", sessionCtx.Err())
			return
		case action, ok := <-heartbeat:
			if !ok {
				heartbeat = nil
				continue
			}
			logger.Printf("session heartbeat: %s", action)
		}
	}
}

func simulateUserActivity(ctx context.Context, heartbeat chan<- string, logger *log.Logger) {
	defer close(heartbeat)
	for _, action := range []string{"search train", "add passengers", "review fare"} {
		select {
		case <-ctx.Done():
			return
		case heartbeat <- action:
			logger.Printf("user action captured: %s", action)
		}
		time.Sleep(1 * time.Second)
	}
	logger.Printf("user went idle; waiting for auto-logout")
}

func seatQuotaDemo(logger *log.Logger) {
	type seatQuota struct {
		sync.Mutex
		tatkal  int
		general int
	}

	quota := &seatQuota{tatkal: 400, general: 1600}
	var wg sync.WaitGroup

	book := func(kind string, seats int) {
		quota.Lock()
		defer quota.Unlock()
		if kind == "tatkal" && quota.tatkal >= seats {
			quota.tatkal -= seats
			return
		}
		if kind == "general" && quota.general >= seats {
			quota.general -= seats
		}
	}

	requests := []struct {
		kind  string
		seats int
	}{
		{"tatkal", 2}, {"tatkal", 1}, {"general", 4}, {"general", 6}, {"general", 3},
		{"tatkal", 3}, {"general", 8}, {"tatkal", 2}, {"general", 5}, {"general", 10},
	}

	for _, req := range requests {
		wg.Add(1)
		go func(r struct {
			kind  string
			seats int
		}) {
			defer wg.Done()
			for i := 0; i < r.seats; i++ {
				book(r.kind, 1)
				time.Sleep(20 * time.Millisecond)
			}
		}(req)
	}

	wg.Wait()
	logger.Printf("seat quota remaining: tatkal=%d general=%d", quota.tatkal, quota.general)
}

func loginFloodSimulation(logger *log.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan int, workerCount*2)
	results := make(chan error, workerCount)
	var processed int64

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for userID := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := authenticateUser(userID); err != nil {
					results <- fmt.Errorf("worker-%d user-%d: %w", id, userID, err)
					return
				}
				total := atomic.AddInt64(&processed, 1)
				if total%200000 == 0 {
					logger.Printf("absorbed %d/%d logins", total, totalConcurrentLogins)
				}
			}
			results <- nil
		}(i)
	}

	go func() {
		defer close(jobs)
		for u := 1; u <= totalConcurrentLogins; u++ {
			select {
			case <-ctx.Done():
				return
			case jobs <- u:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var firstErr error
	for err := range results {
		if err != nil && firstErr == nil {
			firstErr = err
			cancel()
		}
	}

	if firstErr != nil {
		logger.Printf("login flood throttled due to: %v", firstErr)
		return
	}
	logger.Printf("login flood completed: %d users via %d workers", totalConcurrentLogins, workerCount)
}

func authenticateUser(userID int) error {
	if userID%333333 == 0 {
		return errors.New("OTP service saturation")
	}
	return nil
}

func ticketPipeline(logger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	trainID := "12951"
	coachAllocation := make(chan string, 1)

	g.Go(func() error {
		time.Sleep(150 * time.Millisecond)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case coachAllocation <- "B2":
			return nil
		}
	})

	g.Go(func() error {
		return lockSeats(ctx, trainID, 2)
	})

	g.Go(func() error {
		return processUPIPayment(ctx, "txn-8842")
	})

	if err := g.Wait(); err != nil {
		logger.Printf("ticket pipeline aborted: %v", err)
		return
	}

	close(coachAllocation)
	for coach := range coachAllocation {
		logger.Printf("coach assigned: %s", coach)
	}
	logger.Printf("ticket pipeline confirmed for train %s", trainID)
}

func lockSeats(ctx context.Context, trainID string, seats int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(120 * time.Millisecond):
		return nil
	}
}

func processUPIPayment(ctx context.Context, txn string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(200 * time.Millisecond):
		return nil
	}
}

func controlPlane(logger *log.Logger) {
	var g run.Group
	sessionCtx, stopSessions := context.WithCancel(context.Background())
	paymentCtx, stopPayments := context.WithCancel(context.Background())

	g.Add(func() error {
		logger.Printf("session service online")
		<-sessionCtx.Done()
		logger.Printf("session service draining")
		return sessionCtx.Err()
	}, func(err error) {
		logger.Printf("session interrupt: %v", err)
		stopSessions()
	})

	g.Add(func() error {
		logger.Printf("payment switch accepting txns")
		<-paymentCtx.Done()
		logger.Printf("payment switch halting txns")
		return paymentCtx.Err()
	}, func(err error) {
		logger.Printf("payment interrupt: %v", err)
		stopPayments()
	})

	g.Add(func() error {
		time.Sleep(1 * time.Second)
		return errors.New("primary datacenter outage")
	}, func(error) {
		logger.Printf("activating DR drills")
	})

	if err := g.Run(); err != nil {
		logger.Printf("control plane exit: %v", err)
	}
}
