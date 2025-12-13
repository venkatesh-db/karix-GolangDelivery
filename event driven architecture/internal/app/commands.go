package app

import "github.com/venkatesh/order-service/internal/domain"

// PlaceOrder carries user intent to create a new order aggregate.
type PlaceOrder struct {
	OrderID    string
	CustomerID string
	Items      []domain.LineItem
}

// AuthorizePayment is triggered by payment service callback.
type AuthorizePayment struct {
	OrderID   string
	PaymentID string
	Amount    int64
}

// ReserveInventory represents the inventory locking step.
type ReserveInventory struct {
	OrderID       string
	ReservationID string
}

// ShipOrder finalizes fulfillment.
type ShipOrder struct {
	OrderID        string
	TrackingNumber string
	Carrier        string
}

// CancelOrder compensates sagas when downstream failures occur.
type CancelOrder struct {
	OrderID string
	Reason  string
}
