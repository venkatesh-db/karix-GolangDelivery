package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentInitiated PaymentStatus = "initiated"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	SubscriptionID uuid.UUID
	Amount         int64
	Currency       string
	Status         PaymentStatus
	GatewayRef     string
	PaymentDate    time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
