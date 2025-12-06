package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oklog/run"
	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logger := log.New(os.Stdout, "[demo] ", log.LstdFlags|log.Lmicroseconds)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	logger.Println("context-managed goroutines")
	runContextManagedWorkers(ctx, logger)

	logger.Println("channel coordination")
	channelDemo(logger)

	logger.Println("mutex-protected shared state")
	mutexDemo(logger)

	logger.Println("errgroup fail-fast orchestration")
	errgroupDemo(logger)

	logger.Println("oklog/run lifecycle management")
	runGroupDemo(logger)

	logger.Println("Jio/Hotstar-inspired orchestration")
	jioHotstarCaseStudy(logger)
}

func runContextManagedWorkers(parent context.Context, logger *log.Logger) {
	ctx, cancel := context.WithTimeout(parent, 1200*time.Millisecond)
	defer cancel()

	tasks := []time.Duration{300 * time.Millisecond, 700 * time.Millisecond, 1500 * time.Millisecond}
	var wg sync.WaitGroup

	for idx, d := range tasks {
		wg.Add(1)
		go func(id int, delay time.Duration) {
			defer wg.Done()
			if err := fetchWithContext(ctx, id, delay); err != nil {
				logger.Printf("worker %d stopped: %v", id, err)
				return
			}
			logger.Printf("worker %d finished", id)
		}(idx+1, d)
	}

	wg.Wait()
}

func fetchWithContext(ctx context.Context, id int, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return fmt.Errorf("worker %d canceled: %w", id, ctx.Err())
	}
}

func channelDemo(logger *log.Logger) {
	unbuffered := make(chan string)
	go func() { unbuffered <- "handshake complete" }()
	logger.Printf("unbuffered channel message: %s", <-unbuffered)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	buffered := make(chan int, 2)
	go func() {
		defer close(buffered)
		for n := 0; n < 5; n++ {
			select {
			case <-ctx.Done():
				return
			case buffered <- n:
				logger.Printf("produced %d", n)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("channel consumer timeout: %v", ctx.Err())
			return
		case v, ok := <-buffered:
			if !ok {
				logger.Printf("buffered channel drained")
				return
			}
			logger.Printf("consumed %d", v)
		}
	}
}

func mutexDemo(logger *log.Logger) {
	type counter struct {
		sync.Mutex
		value int
	}

	c := &counter{}
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Lock()
				c.value++
				c.Unlock()
			}
		}()
	}

	wg.Wait()
	logger.Printf("counter result %d", c.value)
}

func errgroupDemo(logger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	urls := []string{"https://api.service/a", "https://api.service/b", "https://api.service/c"}
	results := make(map[string]string)
	var mu sync.Mutex

	for _, u := range urls {
		url := u
		g.Go(func() error {
			body, err := fetchEndpoint(ctx, url)
			if err != nil {
				return err
			}
			mu.Lock()
			results[url] = body
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		logger.Printf("errgroup aborted: %v", err)
		return
	}

	for url, body := range results {
		logger.Printf("%s => %s", url, body)
	}
}

func fetchEndpoint(ctx context.Context, url string) (string, error) {
	delay := 200*time.Millisecond + time.Duration(rand.Intn(400))*time.Millisecond
	select {
	case <-time.After(delay):
		if strings.Contains(url, "b") {
			return "", fmt.Errorf("remote error at %s", url)
		}
		return fmt.Sprintf("body-from-%s", url), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// runGroupDemo mirrors how a streaming platform (think Netflix/Hotstar control plane)
// coordinates long-lived processes: one failure triggers a system-wide, ordered shutdown.
func runGroupDemo(logger *log.Logger) {
	var g run.Group
	serviceCtx, cancel := context.WithCancel(context.Background())

	g.Add(func() error {
		logger.Printf("data plane started")
		<-serviceCtx.Done()
		logger.Printf("data plane stopping")
		return serviceCtx.Err()
	}, func(err error) {
		logger.Printf("data plane interrupt: %v", err)
		cancel()
	})

	g.Add(func() error {
		time.Sleep(800 * time.Millisecond)
		return errors.New("cache node lost")
	}, func(error) {
		logger.Printf("cache node cleanup complete")
	})

	if err := g.Run(); err != nil {
		logger.Printf("run group exit: %v", err)
	}
}

type segmentEvent struct {
	session string
	index   int
	bitrate string
}

type cacheMetrics struct {
	sync.Mutex
	hits int
}

func (c *cacheMetrics) add(delta int) {
	c.Lock()
	c.hits += delta
	c.Unlock()
}

func jioHotstarCaseStudy(logger *log.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	logger.Println("JioHotstar: playback session with context & channels")
	if err := playbackSession(ctx, "session-jh-42", logger); err != nil {
		logger.Printf("playback termination: %v", err)
	}

	logger.Println("JioHotstar: cache metrics via buffered channel + mutex")
	cacheUpdatePipeline(ctx, logger)

	logger.Println("JioHotstar: multi-CDN fetch with errgroup")
	if err := cdnFanout(ctx, logger); err != nil {
		logger.Printf("cdn fan-out error: %v", err)
	}

	logger.Println("JioHotstar: platform lifecycle with oklog/run")
	hotstarLifecycle(logger)
}

func playbackSession(ctx context.Context, sessionID string, logger *log.Logger) error {
	sessionCtx, cancel := context.WithTimeout(ctx, 1800*time.Millisecond)
	defer cancel()

	segments := make(chan segmentEvent, 2)
	go prefetchSegments(sessionCtx, sessionID, segments, logger)

	for {
		select {
		case <-sessionCtx.Done():
			return fmt.Errorf("session %s ended: %w", sessionID, sessionCtx.Err())
		case evt, ok := <-segments:
			if !ok {
				logger.Printf("session %s drained buffer", sessionID)
				return nil
			}
			logger.Printf("session %s playing segment %d (%s)", sessionID, evt.index, evt.bitrate)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func prefetchSegments(ctx context.Context, sessionID string, out chan<- segmentEvent, logger *log.Logger) {
	defer close(out)
	bitrates := []string{"480p", "720p", "1080p"}
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			logger.Printf("segment prefetch canceled for %s", sessionID)
			return
		case out <- segmentEvent{session: sessionID, index: i, bitrate: bitrates[rand.Intn(len(bitrates))]}:
			logger.Printf("CDN queued segment %d for %s", i, sessionID)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func cacheUpdatePipeline(ctx context.Context, logger *log.Logger) {
	cache := &cacheMetrics{}
	updates := make(chan int, 3)
	go func() {
		defer close(updates)
		deltas := []int{5, 3, -1, 7}
		for _, delta := range deltas {
			select {
			case <-ctx.Done():
				return
			case updates <- delta:
				logger.Printf("edge cache delta queued: %d", delta)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("metrics pipeline canceled: %v", ctx.Err())
			return
		case delta, ok := <-updates:
			if !ok {
				logger.Printf("JioHotstar cache hits total: %d", cache.hits)
				return
			}
			cache.add(delta)
			logger.Printf("cache hits adjusted by %d => %d", delta, cache.hits)
		}
	}
}

func cdnFanout(ctx context.Context, logger *log.Logger) error {
	g, ctx := errgroup.WithContext(ctx)
	cdns := []string{
		"https://cdn-a.jiohotstar.com/v1/chunk.m4s",
		"https://cdn-b.jiohotstar.com/v1/chunk.m4s",
		"https://cdn-c.jiohotstar.com/v1/chunk.m4s",
	}
	results := make(map[string]string)
	var mu sync.Mutex

	for _, endpoint := range cdns {
		url := endpoint
		g.Go(func() error {
			body, err := fetchCDNSegment(ctx, url)
			if err != nil {
				return err
			}
			mu.Lock()
			results[url] = body
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for url, body := range results {
		logger.Printf("CDN response %s => %s", url, body)
	}
	return nil
}

func fetchCDNSegment(ctx context.Context, endpoint string) (string, error) {
	delay := 150*time.Millisecond + time.Duration(rand.Intn(250))*time.Millisecond
	select {
	case <-time.After(delay):
		if strings.Contains(endpoint, "cdn-b") {
			return "", fmt.Errorf("origin timeout at %s", endpoint)
		}
		return fmt.Sprintf("segment-from-%s", endpoint), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func hotstarLifecycle(logger *log.Logger) {
	var group run.Group
	ingestCtx, stopIngest := context.WithCancel(context.Background())
	playbackCtx, stopPlayback := context.WithCancel(context.Background())

	group.Add(func() error {
		logger.Printf("ingest service online")
		<-ingestCtx.Done()
		logger.Printf("ingest service draining")
		return ingestCtx.Err()
	}, func(err error) {
		logger.Printf("ingest interrupt: %v", err)
		stopIngest()
	})

	group.Add(func() error {
		logger.Printf("playback API serving viewers")
		<-playbackCtx.Done()
		logger.Printf("playback API closed sessions")
		return playbackCtx.Err()
	}, func(err error) {
		logger.Printf("playback interrupt: %v", err)
		stopPlayback()
	})

	group.Add(func() error {
		time.Sleep(900 * time.Millisecond)
		return errors.New("edge cache unavailable")
	}, func(error) {
		logger.Printf("edge cache cleanup done")
	})

	if err := group.Run(); err != nil {
		logger.Printf("hotstar lifecycle exit: %v", err)
	}
}
