package domain

import "time"

// Event represents a domain event emitted by the Order aggregate.
type Event interface {
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}

// BaseEvent contains metadata shared by all events.
type BaseEvent struct {
	Name      string
	EntityID  string
	Timestamp time.Time
}

// EventName returns the event's name.
func (b BaseEvent) EventName() string { return b.Name }

// AggregateID returns the aggregate identifier this event belongs to.
func (b BaseEvent) AggregateID() string { return b.EntityID }

// OccurredAt exposes the timestamp for ordering/replay semantics.
func (b BaseEvent) OccurredAt() time.Time { return b.Timestamp }

// OrderPlaced marks the creation of a new order.
type OrderPlaced struct {
	BaseEvent
	CustomerID string
	Items      []LineItem
	TotalCents int64
}

// PaymentAuthorized indicates the payment service approved the charge.
type PaymentAuthorized struct {
	BaseEvent
	PaymentID string
	Amount    int64
}

// InventoryReserved signals the inventory service locked stock.
type InventoryReserved struct {
	BaseEvent
	ReservationID string
}

// OrderShipped denotes completion of fulfillment.
type OrderShipped struct {
	BaseEvent
	TrackingNumber string
	Carrier        string
}

// OrderCancelled captures compensating workflows when something fails.
type OrderCancelled struct {
	BaseEvent
	Reason string
}
