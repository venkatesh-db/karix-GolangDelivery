package handler

import (
	"encoding/json"
	"net/http"
	"subcription/internal/billing/service"

	"github.com/google/uuid"
)

type BillingHandler struct {

	service *service.BillingService
}

func NewBillingHandler(s *service.BillingService) *BillingHandler {
	return &BillingHandler{
		service: s,
	}
}

type SubscribeRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	PlanID    string `json:"plan_id"`
	Amount    int64 `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

func (h *BillingHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.CreateSubscription(req.UserID, req.PlanID, req.Amount)
	if err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscription created successfully"))
}



