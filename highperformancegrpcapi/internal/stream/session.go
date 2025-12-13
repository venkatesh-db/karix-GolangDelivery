package stream

import (
	"context"
	"sync/atomic"
	"time"

	ridev1 "github.com/example/highperformancegrpcapi/gen/go/proto/ride/v1"
	"github.com/google/uuid"
)

// Session represents a logical connection regardless of the underlying transport.
type Session struct {
	ID         string
	UserID     string
	Transport  string
	outbound   chan *ridev1.ServerEnvelope
	closed     chan struct{}
	cancel     context.CancelFunc
	createdAt  time.Time
	lastSeenMs atomic.Int64
}

// NewSession creates a session with the provided channel capacity.
func NewSession(ctx context.Context, userID, transport string, bufferSize int) (*Session, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s := &Session{
		ID:        uuid.NewString(),
		UserID:    userID,
		Transport: transport,
		outbound:  make(chan *ridev1.ServerEnvelope, bufferSize),
		closed:    make(chan struct{}),
		cancel:    cancel,
		createdAt: time.Now().UTC(),
	}
	s.Touch()
	return s, ctx
}

// Touch updates the last seen timestamp.
func (s *Session) Touch() {
	s.lastSeenMs.Store(time.Now().UTC().UnixMilli())
}

// LastSeen returns when we last heard from the client.
func (s *Session) LastSeen() time.Time {
	return time.UnixMilli(s.lastSeenMs.Load()).UTC()
}

// Outbound returns the channel for pushing envelopes toward the client.
func (s *Session) Outbound() <-chan *ridev1.ServerEnvelope {
	return s.outbound
}

// Enqueue attempts to send the envelope, returning false if the buffer is full.
func (s *Session) Enqueue(msg *ridev1.ServerEnvelope) bool {
	select {
	case s.outbound <- msg:
		return true
	default:
		return false
	}
}

// Close tears down the session.
func (s *Session) Close() {
	select {
	case <-s.closed:
		return
	default:
		close(s.closed)
		close(s.outbound)
		s.cancel()
	}
}

// Done exposes a channel closed when the session terminates.
func (s *Session) Done() <-chan struct{} { return s.closed }

// CreatedAt returns the timestamp when the session was established.
func (s *Session) CreatedAt() time.Time { return s.createdAt }
