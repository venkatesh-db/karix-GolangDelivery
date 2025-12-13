package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	ridepb "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	var target string
	flag.StringVar(&target, "target", "127.0.0.1:7443", "gRPC server address")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Retry dial briefly to avoid race if server is not yet listening
	var conn *grpc.ClientConn
	var err error
	backoff := 200 * time.Millisecond
	for i := 0; i < 10; i++ {
		conn, err = grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
		time.Sleep(backoff)
		if backoff < 2*time.Second {
			backoff *= 2
		}
	}
	if conn == nil || err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := ridepb.NewRideStreamServiceClient(conn)

	md := metadata.Pairs("user-id", "user-123")
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.Connect(ctx)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			srvEnv, err := stream.Recv()
			if err != nil {
				log.Printf("recv done: %v", err)
				return
			}
			switch b := srvEnv.Body.(type) {
			case *ridepb.ServerEnvelope_Ack:
				fmt.Printf("ACK: ok=%v detail=%s\n", b.Ack.Success, b.Ack.Detail)
			case *ridepb.ServerEnvelope_MatchEvent:
				fmt.Printf("MATCH: rider=%s driver=%s eta=%ds surge=%.2f\n", b.MatchEvent.RiderId, b.MatchEvent.DriverId, b.MatchEvent.EtaSeconds, b.MatchEvent.SurgeMultiplier)
			case *ridepb.ServerEnvelope_BroadcastEvent:
				fmt.Printf("BROADCAST: %s len=%d\n", b.BroadcastEvent.Topic, len(b.BroadcastEvent.Payload))
			default:
				fmt.Printf("SERVER: %#v\n", srvEnv)
			}
		}
	}()

	if err := stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_Heartbeat{Heartbeat: &ridepb.Heartbeat{UserId: "user-123", Seq: 1, SentAtUnixMillis: time.Now().UnixMilli()}}}); err != nil {
		log.Fatalf("send heartbeat: %v", err)
	}
	if err := stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_LocationUpdate{LocationUpdate: &ridepb.LocationUpdate{UserId: "user-123", Latitude: 37.7749, Longitude: -122.4194, Bearing: 0, SpeedMps: 0, Sequence: 1, SentAtUnixMillis: time.Now().UnixMilli()}}}); err != nil {
		log.Fatalf("send location: %v", err)
	}
	if err := stream.Send(&ridepb.ClientEnvelope{Body: &ridepb.ClientEnvelope_RideStatusUpdate{RideStatusUpdate: &ridepb.RideStatusUpdate{RideId: "ride-xyz", UserId: "user-123", Status: ridepb.RideStatusUpdate_STATUS_LOOKING, SentAtUnixMillis: time.Now().UnixMilli()}}}); err != nil {
		log.Fatalf("send status: %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}

	if err := stream.CloseSend(); err != nil {
		log.Printf("close send: %v", err)
	}
}
