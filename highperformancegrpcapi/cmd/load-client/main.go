package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	ridepb "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type worker struct {
	id   int
	conn *grpc.ClientConn
	cli  ridepb.RideStreamServiceClient
}

func (w *worker) run(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	defer wg.Done()
	md := metadata.Pairs("user-id", fmt.Sprintf("user-%d", w.id))
	sctx := metadata.NewOutgoingContext(ctx, md)
	stream, err := w.cli.Connect(sctx)
	if err != nil {
		log.Printf("[%d] connect err: %v", w.id, err)
		return
	}

	// Reader goroutine to drain server messages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			srvEnv, err := stream.Recv()
			if err != nil {
				return
			}
			_ = srvEnv // In real load, we could count acks/matches
		}
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	seq := int64(1)
	latBase := 37.7749 + rand.Float64()*0.01
	lonBase := -122.4194 + rand.Float64()*0.01

	for {
		select {
		case <-ctx.Done():
			_ = stream.CloseSend()
			return
		case <-ticker.C:
			// Heartbeat
			_ = stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_Heartbeat{Heartbeat: &ridepb.Heartbeat{UserId: fmt.Sprintf("user-%d", w.id), Seq: seq, SentAtUnixMillis: time.Now().UnixMilli()}}})
			// Location update with slight movement
			lat := latBase + (rand.Float64()-0.5)*0.0005
			lon := lonBase + (rand.Float64()-0.5)*0.0005
			_ = stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_LocationUpdate{LocationUpdate: &ridepb.LocationUpdate{UserId: fmt.Sprintf("user-%d", w.id), Latitude: lat, Longitude: lon, Bearing: 0, SpeedMps: 3.5, Sequence: seq, SentAtUnixMillis: time.Now().UnixMilli()}}})
			// Occasional LOOKING status to trigger matches
			if seq%20 == 0 {
				_ = stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_RideStatusUpdate{RideStatusUpdate: &ridepb.RideStatusUpdate{RideId: fmt.Sprintf("ride-%d", w.id), UserId: fmt.Sprintf("user-%d", w.id), Status: ridepb.RideStatusUpdate_STATUS_LOOKING, SentAtUnixMillis: time.Now().UnixMilli()}}})
			}
			seq++
		}
	}
}

func main() {
	var (
		target   string
		clients  int
		interval time.Duration
		duration time.Duration
	)
	flag.StringVar(&target, "target", "localhost:7443", "gRPC server address")
	flag.IntVar(&clients, "clients", 20, "number of concurrent clients")
	flag.DurationVar(&interval, "interval", 1*time.Second, "send interval per client")
	flag.DurationVar(&duration, "duration", 20*time.Second, "total run duration")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Reuse connections in batches to simulate realistic pooling
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	cli := ridepb.NewRideStreamServiceClient(conn)

	var wg sync.WaitGroup
	wg.Add(clients)
	for i := 0; i < clients; i++ {
		w := &worker{id: i + 1, conn: conn, cli: cli}
		go w.run(ctx, &wg, interval)
	}
	wg.Wait()
}
