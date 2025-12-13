package readmodel

import (
	"context"
	"sort"
	"sync"

	"github.com/venkatesh/order-service/internal/domain"
)

// OrderView is a flattened, query-optimized representation.
type OrderView struct {
	OrderID     string
	CustomerID  string
	Status      domain.OrderStatus
	TotalCents  int64
	LastUpdated int64
}

// OrdersProjection subscribes to all order events and maintains a cache.
type OrdersProjection struct {
	mu      sync.RWMutex
	records map[string]OrderView
}

// NewOrdersProjection initializes the projector map.
func NewOrdersProjection() *OrdersProjection {
	return &OrdersProjection{records: make(map[string]OrderView)}
}

// Handle processes each domain event.
func (p *OrdersProjection) Handle(ctx context.Context, event domain.Event) error {
	_ = ctx
	p.mu.Lock()
	defer p.mu.Unlock()

	view := p.records[event.AggregateID()]
	view.OrderID = event.AggregateID()
	view.LastUpdated = event.OccurredAt().Unix()

	switch e := event.(type) {
	case domain.OrderPlaced:
		view.CustomerID = e.CustomerID
		view.TotalCents = e.TotalCents
		view.Status = domain.OrderStatusConfirmed
	case domain.PaymentAuthorized:
		view.Status = domain.OrderStatusConfirmed
	case domain.InventoryReserved:
		view.Status = domain.OrderStatusReserved
	case domain.OrderShipped:
		view.Status = domain.OrderStatusShipped
	case domain.OrderCancelled:
		view.Status = domain.OrderStatusCancelled
	}

	p.records[view.OrderID] = view
	return nil
}

// List returns all orders sorted by last update desc.
func (p *OrdersProjection) List() []OrderView {
	p.mu.RLock()
	defer p.mu.RUnlock()

	out := make([]OrderView, 0, len(p.records))
	for _, v := range p.records {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastUpdated > out[j].LastUpdated
	})
	return out
}

// Get retrieves a single order by ID.
func (p *OrdersProjection) Get(orderID string) (OrderView, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.records[orderID]
	return v, ok
}
