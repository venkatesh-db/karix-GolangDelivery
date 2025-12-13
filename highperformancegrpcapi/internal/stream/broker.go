package stream

import (
	"context"
	"errors"
	"hash/fnv"
	"sync"
	"time"

	ridev1 "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"go.uber.org/atomic"

	"github.com/example/highperformancegrpcapi/internal/telemetry"
)

// Broker coordinates sessions across shards and handles fan-out semantics.
type Broker struct {
	shards      []*shard
	metrics     *telemetry.Metrics
	maxSessions int
	sessionCnt  atomic.Int64
}

// ErrCapacityReached is returned when the system hit the configured ceiling.
var ErrCapacityReached = errors.New("broker capacity reached")

// NewBroker constructs a Broker with shardCount shards.
func NewBroker(shardCount int, maxSessions int, metrics *telemetry.Metrics) *Broker {
	shards := make([]*shard, shardCount)
	for i := range shards {
		shards[i] = &shard{sessions: make(map[string]map[string]*Session)}
	}
	return &Broker{shards: shards, metrics: metrics, maxSessions: maxSessions}
}

// Register allocates a session for the provided user.
func (b *Broker) Register(ctx context.Context, userID, transport string, bufferSize int) (*Session, context.Context, error) {
	if int(b.sessionCnt.Load()) >= b.maxSessions {
		return nil, nil, ErrCapacityReached
	}

	s, ctxWithCancel := NewSession(ctx, userID, transport, bufferSize)
	sh := b.pick(userID)
	sh.attach(s)

	b.sessionCnt.Inc()
	if b.metrics != nil {
		b.metrics.ActiveSessions.Inc()
	}

	return s, ctxWithCancel, nil
}

// Detach removes the session and publishes metrics.
func (b *Broker) Detach(s *Session) {
	sh := b.pick(s.UserID)
	if sh.detach(s) {
		b.sessionCnt.Dec()
		if b.metrics != nil {
			b.metrics.ActiveSessions.Dec()
			b.metrics.SessionDuration.Observe(timeSince(s.CreatedAt()))
		}
	}
}

// Send sends to all sessions for the given user.
func (b *Broker) Send(userID string, msg *ridev1.ServerEnvelope) int {
	sh := b.pick(userID)
	return sh.broadcastUser(userID, msg, b.metrics)
}

// Broadcast iterates through every session and invokes predicate to decide delivery.
func (b *Broker) Broadcast(predicate func(*Session) bool, msgFactory func(*Session) *ridev1.ServerEnvelope) int {
	delivered := 0
	for _, sh := range b.shards {
		delivered += sh.broadcast(predicate, msgFactory, b.metrics)
	}
	return delivered
}

func (b *Broker) pick(userID string) *shard {
	h := fnv.New64a()
	_, _ = h.Write([]byte(userID))
	idx := int(h.Sum64()) % len(b.shards)
	if idx < 0 {
		idx = -idx
	}
	return b.shards[idx]
}

type shard struct {
	mu       sync.RWMutex
	sessions map[string]map[string]*Session // userID -> sessionID -> Session
}

func (s *shard) attach(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[session.UserID]; !ok {
		s.sessions[session.UserID] = make(map[string]*Session)
	}
	s.sessions[session.UserID][session.ID] = session
}

func (s *shard) detach(session *Session) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	userSessions, ok := s.sessions[session.UserID]
	if !ok {
		return false
	}
	if _, ok := userSessions[session.ID]; !ok {
		return false
	}
	delete(userSessions, session.ID)
	if len(userSessions) == 0 {
		delete(s.sessions, session.UserID)
	}
	return true
}

func (s *shard) broadcast(predicate func(*Session) bool, msgFactory func(*Session) *ridev1.ServerEnvelope, metrics *telemetry.Metrics) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, userSessions := range s.sessions {
		for _, session := range userSessions {
			if predicate != nil && !predicate(session) {
				continue
			}
			msg := msgFactory(session)
			if session.Enqueue(msg) {
				count++
				if metrics != nil {
					metrics.EgressMessages.WithLabelValues(msgBodyLabel(msg), session.Transport).Inc()
				}
			} else if metrics != nil {
				metrics.DroppedMessages.WithLabelValues("queue_full").Inc()
			}
		}
	}
	return count
}

func (s *shard) broadcastUser(userID string, msg *ridev1.ServerEnvelope, metrics *telemetry.Metrics) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userSessions, ok := s.sessions[userID]
	if !ok {
		return 0
	}
	delivered := 0
	for _, session := range userSessions {
		if session.Enqueue(msg) {
			delivered++
			if metrics != nil {
				metrics.EgressMessages.WithLabelValues(msgBodyLabel(msg), session.Transport).Inc()
			}
		} else if metrics != nil {
			metrics.DroppedMessages.WithLabelValues("queue_full").Inc()
		}
	}
	return delivered
}

func msgBodyLabel(msg *ridev1.ServerEnvelope) string {
	switch msg.GetBody().(type) {
	case *ridev1.ServerEnvelope_MatchEvent:
		return "match_event"
	case *ridev1.ServerEnvelope_Ack:
		return "ack"
	case *ridev1.ServerEnvelope_BroadcastEvent:
		return "broadcast"
	default:
		return "unknown"
	}
}

func timeSince(t time.Time) float64 {
	return time.Since(t).Seconds()
}
