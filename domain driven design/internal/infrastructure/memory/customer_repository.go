package memory

import (
	"context"
	"sync"

	"github.com/helrachar/banking/internal/domain"
)

// CustomerRepository is an in-memory repository useful for local demos.
type CustomerRepository struct {
	mu        sync.RWMutex
	customers map[domain.CustomerID]*domain.Customer
}

// NewCustomerRepository constructs the in-memory map store.
func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{customers: make(map[domain.CustomerID]*domain.Customer)}
}

// Save persists or replaces the aggregate.
func (r *CustomerRepository) Save(_ context.Context, customer *domain.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.customers[customer.ID()] = customer
	return nil
}

// GetByID fetches an aggregate or yields ErrNotFound.
func (r *CustomerRepository) GetByID(_ context.Context, id domain.CustomerID) (*domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	customer, ok := r.customers[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return customer, nil
}

// List returns a snapshot slice of all aggregates.
func (r *CustomerRepository) List(_ context.Context) ([]*domain.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]*domain.Customer, 0, len(r.customers))
	for _, c := range r.customers {
		results = append(results, c)
	}
	return results, nil
}
