package server

import (
	"context"
	"io"

	ridev1 "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"github.com/example/highperformancegrpcapi/internal/matching"
	"github.com/example/highperformancegrpcapi/internal/stream"
	"github.com/example/highperformancegrpcapi/internal/telemetry"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

// RideStreamHandler implements the gRPC RideStreamService.
type RideStreamHandler struct {
	ridev1.UnimplementedRideStreamServiceServer

	Broker  *stream.Broker
	Engine  *matching.Engine
	Metrics *telemetry.Metrics
	Log     *zap.Logger
	BufSize int
}

// Connect establishes a bidirectional stream for session level messaging.
func (h *RideStreamHandler) Connect(stream ridev1.RideStreamService_ConnectServer) error {
	ctx := stream.Context()
	userID := userIDFromMetadata(ctx)

	logger := h.Log
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.With(zap.String("user_id", userID))

	session, sessionCtx, err := h.Broker.Register(ctx, userID, "grpc", h.BufSize)
	if err != nil {
		return err
	}
	defer func() {
		h.Broker.Detach(session)
		session.Close()
	}()

	grp, grpCtx := errgroup.WithContext(sessionCtx)

	grp.Go(func() error {
		for {
			env, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				logger.Warn("ingress stream closed", zap.Error(err))
				return err
			}

			session.Touch()
			if h.Metrics != nil {
				h.Metrics.IngressMessages.WithLabelValues(bodyLabel(env), "grpc").Inc()
			}

			h.Engine.Submit(session, env)
		}
	})

	grp.Go(func() error {
		for {
			select {
			case <-grpCtx.Done():
				return grpCtx.Err()
			case msg, ok := <-session.Outbound():
				if !ok {
					return nil
				}
				if err := stream.Send(msg); err != nil {
					logger.Warn("egress stream closed", zap.Error(err))
					return err
				}
			}
		}
	})

	return grp.Wait()
}

func userIDFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.NewString()
	}
	if v := md.Get("user-id"); len(v) > 0 && v[0] != "" {
		return v[0]
	}
	return uuid.NewString()
}

func bodyLabel(env *ridev1.ClientEnvelope) string {
	switch env.GetBody().(type) {
	case *ridev1.ClientEnvelope_LocationUpdate:
		return "location"
	case *ridev1.ClientEnvelope_RideStatusUpdate:
		return "status"
	case *ridev1.ClientEnvelope_Heartbeat:
		return "heartbeat"
	default:
		return "unknown"
	}
}
