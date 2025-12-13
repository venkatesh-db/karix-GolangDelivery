package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/venkatesh/order-service/internal/domain"
)

// InMemoryStore is a thread-safe event store useful for tests/demos.
type InMemoryStore struct {
	mu     sync.RWMutex
	events map[string][]domain.Event
}

// NewInMemoryStore boots an empty store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{events: make(map[string][]domain.Event)}
}

// Load rehydrates the aggregate's event stream.
func (s *InMemoryStore) Load(_ context.Context, aggregateID string) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stream := s.events[aggregateID]
	out := make([]domain.Event, len(stream))
	copy(out, stream)
	return out, nil
}

// Append enforces optimistic concurrency before adding events.
func (s *InMemoryStore) Append(_ context.Context, aggregateID string, expectedVersion int, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current := len(s.events[aggregateID])
	if current != expectedVersion {
		return fmt.Errorf("concurrency conflict for aggregate %s: expected %d got %d", aggregateID, expectedVersion, current)
	}
	s.events[aggregateID] = append(s.events[aggregateID], events...)
	return nil
}
