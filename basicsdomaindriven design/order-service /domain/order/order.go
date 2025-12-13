package order

import "errors"

type Order struct {
	ID     string
	Amount float64
	Status string
}

func NewOrder(id string, amount float64) (*Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	if amount <= 0 {
		return nil, errors.New("order amount must be greater than zero")
	}
	return &Order{
		ID:     id,
		Amount: amount,
		Status: "Created",
	}, nil
}

func (o *Order) Pay() error {

	if o.Status != "Created" {
		return errors.New("only orders in 'Created' status can be paid")
	}
	o.Status = "Paid"
	return nil
}
