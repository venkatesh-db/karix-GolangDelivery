package app

import (
	"context"
	"fmt"

	"github.com/venkatesh/order-service/internal/domain"
)

// EventStore persists and rehydrates aggregates via event sourcing.
type EventStore interface {
	Load(ctx context.Context, aggregateID string) ([]domain.Event, error)
	Append(ctx context.Context, aggregateID string, expectedVersion int, events []domain.Event) error
}

// Publisher fan-outs domain events to interested projections/sagas.
type Publisher interface {
	Publish(ctx context.Context, events []domain.Event) error
}

// OrderService wires command handlers to persistence + messaging.
type OrderService struct {
	store     EventStore
	publisher Publisher
}

// NewOrderService composes an application service instance.
func NewOrderService(store EventStore, publisher Publisher) *OrderService {
	return &OrderService{store: store, publisher: publisher}
}

// HandlePlaceOrder executes the aggregate logic for PlaceOrder.
func (s *OrderService) HandlePlaceOrder(ctx context.Context, cmd PlaceOrder) error {
	if cmd.OrderID == "" {
		return fmt.Errorf("order id is required")
	}
	order := domain.NewOrder(cmd.OrderID)
	expected := order.Version()
	events, err := order.HandlePlaceOrder(cmd.CustomerID, cmd.Items)
	if err != nil {
		return err
	}
	return s.persistAndPublish(ctx, order, expected, events)
}

// HandleAuthorizePayment rehydrates aggregate then delegates to domain logic.
func (s *OrderService) HandleAuthorizePayment(ctx context.Context, cmd AuthorizePayment) error {
	order, err := s.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	expected := order.Version()
	events, err := order.HandleAuthorizePayment(cmd.PaymentID, cmd.Amount)
	if err != nil {
		return err
	}
	return s.persistAndPublish(ctx, order, expected, events)
}

// HandleReserveInventory ensures saga progression after payment.
func (s *OrderService) HandleReserveInventory(ctx context.Context, cmd ReserveInventory) error {
	order, err := s.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	expected := order.Version()
	events, err := order.HandleReserveInventory(cmd.ReservationID)
	if err != nil {
		return err
	}
	return s.persistAndPublish(ctx, order, expected, events)
}

// HandleShipOrder finalizes lifecycle.
func (s *OrderService) HandleShipOrder(ctx context.Context, cmd ShipOrder) error {
	order, err := s.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	expected := order.Version()
	events, err := order.HandleShipOrder(cmd.TrackingNumber, cmd.Carrier)
	if err != nil {
		return err
	}
	return s.persistAndPublish(ctx, order, expected, events)
}

// HandleCancelOrder compensates when downstream services fail.
func (s *OrderService) HandleCancelOrder(ctx context.Context, cmd CancelOrder) error {
	order, err := s.loadOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	expected := order.Version()
	events, err := order.HandleCancel(cmd.Reason)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	return s.persistAndPublish(ctx, order, expected, events)
}

func (s *OrderService) loadOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order id is required")
	}
	events, err := s.store.Load(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order := domain.NewOrder(orderID)
	for _, evt := range events {
		order.Apply(evt)
	}
	return order, nil
}

func (s *OrderService) persistAndPublish(ctx context.Context, order *domain.Order, expected int, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}
	if err := s.store.Append(ctx, order.ID(), expected, events); err != nil {
		return err
	}
	for _, evt := range events {
		order.Apply(evt)
	}
	if err := s.publisher.Publish(ctx, events); err != nil {
		return err
	}
	return nil
}
