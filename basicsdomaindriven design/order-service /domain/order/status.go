package order


type Status string

const (
	StatusCreated Status = "Created"
	StatusPaid    Status = "Paid"

	StatusShipped Status = "Shipped"
	StatusDelivered Status = "Delivered"
	StatusCancelled Status = "Cancelled"
)


