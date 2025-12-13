package application


import (
	"climax/domain/order"
	"errors"
)

type OrderService struct {
	repo order.Repository
}

func NewOrderService(repo order.Repository) *OrderService {
	return &OrderService{repo: repo}
}


func (s *OrderService) PlaceOrder(id string, amount float64) (*order.Order, error) {

	ord, err := order.NewOrder(id, amount)
	if err != nil {
		return nil, errors.New("failed to create order: " + err.Error())
	}
	
	err = s.repo.Save(ord)
	if err != nil {
		return nil, errors.New("failed to save order: " + err.Error())
	}
	return ord, nil
}

