package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	StatusPending   SubscriptionStatus = "pending"
	StatusActive    SubscriptionStatus = "active"
	StatusCancelled SubscriptionStatus = "canceled"
	StatusExpired   SubscriptionStatus = "expired"
)

type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PlanID    string
	Amount    int64
	Status    SubscriptionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
