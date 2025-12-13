package subscriptions

import "time"

// Plan describes a sellable subscription plan persisted in Postgres.
type Plan struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PriceCents    int       `json:"price_cents"`
	BillingPeriod string    `json:"billing_period"`
	MaxSeats      int       `json:"max_seats"`
	CreatedAt     time.Time `json:"created_at"`
}

// TenantSubscription captures the active plan for a tenant.
type TenantSubscription struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	PlanID           string    `json:"plan_id"`
	Status           string    `json:"status"`
	Seats            int       `json:"seats"`
	ActivatedAt      time.Time `json:"activated_at"`
	CurrentPeriodEnd time.Time `json:"current_period_end"`
}

// ActivateInput describes the payload to activate or switch a plan.
type ActivateInput struct {
	TenantID string `json:"tenant_id"`
	PlanID   string `json:"plan_id"`
	Seats    int    `json:"seats"`
}
