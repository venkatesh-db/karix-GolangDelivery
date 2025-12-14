package main

import (
	"log"
	"net/http"
	"subcription/internal/billing/handler"
	"subcription/internal/billing/repository"
	"subcription/internal/billing/service"
	"subcription/internal/platform/payment"
)

/*

curl -X POST http://localhost:8080/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "plan_id": "basic_plan",
    "amount": 999,
    "payment_method": "card"
  }'


*/

func main() {

	// All dependencies would be initialized here (e.g., database connections, config loading, etc.)

	var repo repository.BillingRepository

	// Correctly initialize the repository
	repo = repository.NewInMemoryBillingRepository()

	gateway := &payment.StripeMock{} // Assume paymentGateway is defined elsewhere
	// Correctly initialize the billing service
	billingService := service.NewBillingService(repo, gateway)
	billingHandler := handler.NewBillingHandler(billingService)

	mux := http.NewServeMux()
	mux.HandleFunc("/subscribe", billingHandler.Subscribe)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

}
