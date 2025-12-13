package main

import (
	"climax/application"
	repository "climax/infrastructure/repository"
	"fmt"
)

func main() {

	repo := repository.NewMemoryOrderRepo()

	service := application.NewOrderService(repo)

	service.PlaceOrder("ORD-1", 100.0)
	fmt.Println("Order ORD-1 placed.")

	order, err := service.PlaceOrder("ORD-1", 100.0)
	if err != nil {
		fmt.Println("Error retrieving order:", err)
		return
	}
	fmt.Printf("Retrieved Order: ID=%s, Amount=%.2f\n", order.ID, order.Amount)

}
