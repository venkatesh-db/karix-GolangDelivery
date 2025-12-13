package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/helrachar/banking/internal/domain"
)

// Service orchestrates customer use-cases within the application layer.
type Service struct {
	repo domain.CustomerRepository
	now  func() time.Time
}

// NewService wires dependencies for the customer application service.
func NewService(repo domain.CustomerRepository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Now().UTC() }}
}

// RegisterCustomer validates and persists a new customer aggregate.
func (s *Service) RegisterCustomer(ctx context.Context, id domain.CustomerID, fullName, email, pan string) (*domain.Customer, error) {
	customer, err := domain.NewCustomer(id, fullName, email, pan, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, customer); err != nil {
		return nil, fmt.Errorf("save customer: %w", err)
	}
	return customer, nil
}

// FetchCustomer returns an aggregate by id.
func (s *Service) FetchCustomer(ctx context.Context, id domain.CustomerID) (*domain.Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return customer, nil
}

// ListCustomers returns all aggregates (demo scope).
func (s *Service) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	return customers, nil
}
