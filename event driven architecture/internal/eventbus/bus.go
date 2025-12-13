package eventbus

import (
	"context"
	"sync"

	"github.com/venkatesh/order-service/internal/domain"
)

// Handler consumes a domain event.
type Handler func(ctx context.Context, event domain.Event) error

// Bus is a lightweight in-memory event dispatcher suited for demos/tests.
type Bus struct {
	mu           sync.RWMutex
	subscribers  map[string][]Handler
	asyncWorkers int
}

// WildcardEvent routes every event to the handler.
const WildcardEvent = "*"

// New creates an event bus with the provided async parallelism.
func New(asyncWorkers int) *Bus {
	if asyncWorkers <= 0 {
		asyncWorkers = 1
	}
	return &Bus{
		subscribers:  make(map[string][]Handler),
		asyncWorkers: asyncWorkers,
	}
}

// Subscribe registers a handler for a given event name or WildcardEvent.
func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventName] = append(b.subscribers[eventName], handler)
}

// Publish fan-outs every event to interested subscribers asynchronously.
func (b *Bus) Publish(ctx context.Context, events []domain.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var wg sync.WaitGroup
	sem := make(chan struct{}, b.asyncWorkers)

	dispatch := func(handler Handler, evt domain.Event) {
		defer wg.Done()
		if err := handler(ctx, evt); err != nil {
			// Errors are intentionally swallowed for brevity; production code
			// should capture metrics and send failures to a DLQ.
		}
		<-sem
	}

	for _, evt := range events {
		handlers := append([]Handler{}, b.subscribers[evt.EventName()]...)
		handlers = append(handlers, b.subscribers[WildcardEvent]...)
		for _, handler := range handlers {
			wg.Add(1)
			sem <- struct{}{}
			go dispatch(handler, evt)
		}
	}

	wg.Wait()
	return nil
}
