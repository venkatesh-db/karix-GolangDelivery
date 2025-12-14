package service

import (
	"subcription/internal/billing/domain"
	"subcription/internal/billing/repository"
	"subcription/internal/platform/payment"
	"time"

	"github.com/google/uuid"
)

type BillingService struct {
	repo    repository.BillingRepository
	gateway payment.Gatway
}

func NewBillingService(repo repository.BillingRepository, gateway payment.Gatway) *BillingService {
	return &BillingService{
		repo:    repo,
		gateway: gateway,
	}
}

func (s *BillingService) CreateSubscription(userID uuid.UUID, planID string, amount int64) error {

	subscription := &domain.Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		PlanID:    planID,
		Amount:    amount,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateSubscription(subscription.UserID, subscription.PlanID, subscription.Amount); err != nil {
		return err
	}

	resp, err := s.gateway.Charge(amount, "usd", "src_mocked_12345", "Subscription Charge")

	if err != nil || !resp.Success {
		return err
	}

	payment := &domain.Payment{
		ID:             uuid.New(),
		SubscriptionID: subscription.ID,
		Amount:         amount,
		Currency:       "INR",
		Status:         "completed",
		CreatedAt:      time.Now(),
	}

	if err := s.repo.CreatePayment(payment); err != nil {
		return err
	}

	return nil

}
