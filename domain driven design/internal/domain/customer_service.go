package domain

import (
	"context"
	"fmt"
	"time"
)

// CustomerRepository abstracts the persistence for the aggregate.
type CustomerRepository interface {
	Save(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id CustomerID) (*Customer, error)
	List(ctx context.Context) ([]*Customer, error)
}

// CustomerService orchestrates simple use-cases around the aggregate.
type CustomerService struct {
	repo CustomerRepository
}

// NewCustomerService wires the repository into the service boundary.
func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// RegisterCustomer creates the customer aggregate and persists it.
func (s *CustomerService) RegisterCustomer(ctx context.Context, id CustomerID, fullName, email, pan string) (*Customer, error) {
	customer, err := NewCustomer(id, fullName, email, pan, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, customer); err != nil {
		return nil, fmt.Errorf("save customer: %w", err)
	}
	return customer, nil
}

// FetchCustomer returns an aggregate by id.
func (s *CustomerService) FetchCustomer(ctx context.Context, id CustomerID) (*Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return customer, nil
}

// ListCustomers is a simple query operation for demo purposes.
func (s *CustomerService) ListCustomers(ctx context.Context) ([]*Customer, error) {
	customers, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	return customers, nil
}
