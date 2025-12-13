package domain

import (
	"errors"
	"fmt"
	"time"
)

// OrderStatus captures the lifecycle of an order aggregate.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusReserved  OrderStatus = "reserved"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// LineItem represents a product requested by the shopper.
type LineItem struct {
	SKU            string
	Quantity       int
	UnitPriceCents int64
}

// Order is the aggregate root responsible for enforcing invariants.
type Order struct {
	id                string
	customerID        string
	items             []LineItem
	status            OrderStatus
	version           int
	totalCents        int64
	paymentAuthorized bool
	inventoryReserved bool
}

// Domain errors.
var (
	ErrInvalidLineItems   = errors.New("order must contain at least one valid line item")
	ErrPaymentMismatch    = errors.New("payment amount does not match order total")
	ErrInvalidTransition  = errors.New("invalid state transition for order")
	ErrAlreadyInitialized = errors.New("order already initialized")
)

// NewOrder builds an empty aggregate; callers must apply events to hydrate it.
func NewOrder(id string) *Order {
	return &Order{
		id:     id,
		status: OrderStatusPending,
	}
}

// ID returns the aggregate identifier.
func (o *Order) ID() string { return o.id }

// Version exposes the optimistic concurrency token.
func (o *Order) Version() int { return o.version }

// Status returns the lifecycle state.
func (o *Order) Status() OrderStatus { return o.status }

// HandlePlaceOrder validates command intent and returns resulting events.
func (o *Order) HandlePlaceOrder(customerID string, items []LineItem) ([]Event, error) {
	if o.version > 0 {
		return nil, ErrAlreadyInitialized
	}
	if len(items) == 0 {
		return nil, ErrInvalidLineItems
	}
	var total int64
	for _, item := range items {
		if item.SKU == "" || item.Quantity <= 0 || item.UnitPriceCents <= 0 {
			return nil, ErrInvalidLineItems
		}
		total += int64(item.Quantity) * item.UnitPriceCents
	}

	evt := OrderPlaced{
		BaseEvent: BaseEvent{
			Name:      "OrderPlaced",
			EntityID:  o.id,
			Timestamp: time.Now().UTC(),
		},
		CustomerID: customerID,
		Items:      items,
		TotalCents: total,
	}

	return []Event{evt}, nil
}

// HandleAuthorizePayment confirms funds before inventory reservation.
func (o *Order) HandleAuthorizePayment(paymentID string, amount int64) ([]Event, error) {
	if o.status == OrderStatusCancelled {
		return nil, ErrInvalidTransition
	}
	if amount != o.totalCents {
		return nil, ErrPaymentMismatch
	}
	if o.paymentAuthorized {
		return nil, fmt.Errorf("payment already authorized for order %s", o.id)
	}

	evt := PaymentAuthorized{
		BaseEvent: BaseEvent{
			Name:      "PaymentAuthorized",
			EntityID:  o.id,
			Timestamp: time.Now().UTC(),
		},
		PaymentID: paymentID,
		Amount:    amount,
	}
	return []Event{evt}, nil
}

// HandleReserveInventory emits InventoryReserved once payment is secured.
func (o *Order) HandleReserveInventory(reservationID string) ([]Event, error) {
	if !o.paymentAuthorized || o.status == OrderStatusCancelled {
		return nil, ErrInvalidTransition
	}
	if o.inventoryReserved {
		return nil, fmt.Errorf("inventory already reserved for order %s", o.id)
	}
	evt := InventoryReserved{
		BaseEvent: BaseEvent{
			Name:      "InventoryReserved",
			EntityID:  o.id,
			Timestamp: time.Now().UTC(),
		},
		ReservationID: reservationID,
	}
	return []Event{evt}, nil
}

// HandleShipOrder finalizes fulfillment after reservation.
func (o *Order) HandleShipOrder(tracking, carrier string) ([]Event, error) {
	if !o.inventoryReserved || o.status == OrderStatusCancelled {
		return nil, ErrInvalidTransition
	}
	if o.status == OrderStatusShipped {
		return nil, fmt.Errorf("order %s already shipped", o.id)
	}
	evt := OrderShipped{
		BaseEvent: BaseEvent{
			Name:      "OrderShipped",
			EntityID:  o.id,
			Timestamp: time.Now().UTC(),
		},
		TrackingNumber: tracking,
		Carrier:        carrier,
	}
	return []Event{evt}, nil
}

// HandleCancel ensures compensating action can be triggered.
func (o *Order) HandleCancel(reason string) ([]Event, error) {
	if o.status == OrderStatusCancelled {
		return nil, nil
	}
	evt := OrderCancelled{
		BaseEvent: BaseEvent{
			Name:      "OrderCancelled",
			EntityID:  o.id,
			Timestamp: time.Now().UTC(),
		},
		Reason: reason,
	}
	return []Event{evt}, nil
}

// Apply mutates aggregate state from an event during command handling or replay.
func (o *Order) Apply(event Event) {
	switch e := event.(type) {
	case OrderPlaced:
		o.customerID = e.CustomerID
		o.items = append([]LineItem(nil), e.Items...)
		o.totalCents = e.TotalCents
		o.status = OrderStatusConfirmed
	case PaymentAuthorized:
		o.paymentAuthorized = true
		o.status = OrderStatusConfirmed
	case InventoryReserved:
		o.inventoryReserved = true
		o.status = OrderStatusReserved
	case OrderShipped:
		o.status = OrderStatusShipped
	case OrderCancelled:
		o.status = OrderStatusCancelled
	}
	o.version++
}
