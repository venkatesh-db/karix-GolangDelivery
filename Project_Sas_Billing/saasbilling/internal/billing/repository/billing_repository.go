package repository

import (
	"subcription/internal/billing/domain"
	"time"

	"github.com/google/uuid"
)

type BillingRepository interface {
	CreateSubscription(userID uuid.UUID, planID string, amount int64) error
	CreatePayment(payment *domain.Payment) error

	//CreateSubscription(sub *domain.Subsciption) error
	//UpdateSubscription(id string , status domain.SubscriptionStatus) error

	//UpdatePaymentStatus(id string, status domain.PaymentStatus,ref string) error

}

// NewInMemoryBillingRepository creates an in-memory implementation of BillingRepository
func NewInMemoryBillingRepository() BillingRepository {
	return &inMemoryBillingRepository{
		subscriptions: make(map[uuid.UUID]*domain.Subscription),
		payments:      make(map[uuid.UUID]*domain.Payment),
	}
}

type inMemoryBillingRepository struct {
	subscriptions map[uuid.UUID]*domain.Subscription
	payments      map[uuid.UUID]*domain.Payment
}

func (r *inMemoryBillingRepository) CreateSubscription(userID uuid.UUID, planID string, amount int64) error {
	sub := &domain.Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		PlanID:    planID,
		Amount:    amount,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	r.subscriptions[sub.ID] = sub
	return nil
}

func (r *inMemoryBillingRepository) CreatePayment(payment *domain.Payment) error {
	r.payments[payment.ID] = payment
	return nil
}
