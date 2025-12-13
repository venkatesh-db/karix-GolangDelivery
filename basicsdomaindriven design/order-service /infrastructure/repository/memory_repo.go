package repository

import (
	"climax/domain/order"
	"errors"
)

type MemoryOrderRepo struct {
	orders map[string]*order.Order
}

func NewMemoryOrderRepo() *MemoryOrderRepo {
	return &MemoryOrderRepo{
		orders: make(map[string]*order.Order),
	}
}

func (r *MemoryOrderRepo) Save(ord *order.Order) error {
	if ord == nil {
		return errors.New("order is nil")
	}
	r.orders[ord.ID] = ord
	return nil
}

func (r *MemoryOrderRepo) FindByID(id string) (*order.Order, error) {
	ord, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}
	return ord, nil
}
