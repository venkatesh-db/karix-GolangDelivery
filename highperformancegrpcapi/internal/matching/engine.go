package matching

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"time"

	ridev1 "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"github.com/example/highperformancegrpcapi/internal/stream"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Engine is a toy stand-in for the dispatch system that would normally live outside this service.
type Engine struct {
	broker  *stream.Broker
	log     *zap.Logger
	workers int
	queue   chan sessionEnvelope
}

type sessionEnvelope struct {
	session *stream.Session
	env     *ridev1.ClientEnvelope
}

// New creates a new Engine with the provided worker count.
func New(broker *stream.Broker, log *zap.Logger, workers int) *Engine {
	return &Engine{
		broker:  broker,
		log:     log.Named("matching"),
		workers: workers,
		queue:   make(chan sessionEnvelope, 64_000),
	}
}

// Start launches background workers.
func (e *Engine) Start(ctx context.Context) {
	for i := 0; i < e.workers; i++ {
		go e.worker(ctx, i)
	}
}

// Submit hands an envelope to the engine.
func (e *Engine) Submit(session *stream.Session, env *ridev1.ClientEnvelope) {
	select {
	case e.queue <- sessionEnvelope{session: session, env: env}:
	default:
		if e.log != nil {
			e.log.Warn("ingress queue saturated", zap.String("user", session.UserID))
		}
	}
}

func (e *Engine) worker(ctx context.Context, id int) {
	logger := e.log.With(zap.Int("worker", id))
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-e.queue:
			if evt.env == nil {
				continue
			}
			e.route(logger, evt.session, evt.env)
		}
	}
}

func (e *Engine) route(logger *zap.Logger, session *stream.Session, env *ridev1.ClientEnvelope) {
	if hb := env.GetHeartbeat(); hb != nil {
		session.Touch()
		return
	}

	switch body := env.GetBody().(type) {
	case *ridev1.ClientEnvelope_LocationUpdate:
		e.handleLocation(session, body.LocationUpdate)
	case *ridev1.ClientEnvelope_RideStatusUpdate:
		e.handleStatus(session, body.RideStatusUpdate)
	default:
		logger.Debug("discarding envelope", zap.String("session", session.ID))
	}
}

func (e *Engine) handleLocation(session *stream.Session, update *ridev1.LocationUpdate) {
	zone := zoneKey(update.Latitude, update.Longitude)
	e.broker.Broadcast(func(target *stream.Session) bool {
		return target.UserID != session.UserID && hashString(target.UserID)%32 == hashString(zone)%32
	}, func(target *stream.Session) *ridev1.ServerEnvelope {
		return &ridev1.ServerEnvelope{
			CorrelationId: uuid.NewString(),
			LamportTime:   envLamport(update.Sequence),
			Body: &ridev1.ServerEnvelope_BroadcastEvent{BroadcastEvent: &ridev1.BroadcastEvent{
				Topic:   zone,
				Payload: marshalLocation(update),
			}},
		}
	})
}

func (e *Engine) handleStatus(session *stream.Session, status *ridev1.RideStatusUpdate) {
	if status.Status == ridev1.RideStatusUpdate_STATUS_LOOKING {
		e.queueMatch(status.RideId, status.UserId)
	}

	ack := &ridev1.ServerEnvelope{
		CorrelationId: status.RideId,
		LamportTime:   envLamport(status.SentAtUnixMillis),
		Body: &ridev1.ServerEnvelope_Ack{Ack: &ridev1.Ack{
			CorrelationId: status.RideId,
			Success:       true,
			Detail:        "status applied",
		}},
	}
	e.broker.Send(session.UserID, ack)
}

func (e *Engine) queueMatch(rideID, riderID string) {
	driverID := fmt.Sprintf("drv-%06d", rand.Intn(900000))
	match := &ridev1.ServerEnvelope{
		CorrelationId: rideID,
		LamportTime:   time.Now().UnixNano(),
		Body: &ridev1.ServerEnvelope_MatchEvent{MatchEvent: &ridev1.MatchEvent{
			RideId:          rideID,
			DriverId:        driverID,
			RiderId:         riderID,
			VehiclePlate:    fmt.Sprintf("TN-%02d-%04d", rand.Intn(99), rand.Intn(9999)),
			EtaSeconds:      int64(60 + rand.Intn(300)),
			SurgeMultiplier: 1 + rand.Float64()*0.5,
		}},
	}
	e.broker.Send(riderID, match)
}

func marshalLocation(update *ridev1.LocationUpdate) []byte {
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint64(buf[0:8], math.Float64bits(update.Latitude))
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(update.Longitude))
	binary.LittleEndian.PutUint64(buf[16:24], math.Float64bits(update.Bearing))
	binary.LittleEndian.PutUint64(buf[24:32], math.Float64bits(update.SpeedMps))
	return buf
}

func zoneKey(lat, lng float64) string {
	return fmt.Sprintf("%d:%d", int(lat*100), int(lng*100))
}

func hashString(v string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(v))
	return h.Sum64()
}

func envLamport(seed int64) int64 {
	if seed == 0 {
		return time.Now().UnixNano()
	}
	return seed
}
