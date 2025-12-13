package order

type Repository interface {
	Save(ord *Order) error
	FindByID(id string) (*Order, error)
}
