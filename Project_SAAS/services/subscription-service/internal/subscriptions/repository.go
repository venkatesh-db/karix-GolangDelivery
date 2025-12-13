package subscriptions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persists plans and tenant subscriptions.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var (
	ErrPlanNotFound         = errors.New("plan not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

func (r *Repository) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, description, price_cents, billing_period, max_seats, created_at FROM plans ORDER BY price_cents ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PriceCents, &p.BillingPeriod, &p.MaxSeats, &p.CreatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *Repository) GetPlan(ctx context.Context, planID string) (Plan, error) {
	var p Plan
	err := r.pool.QueryRow(ctx, `SELECT id, name, description, price_cents, billing_period, max_seats, created_at FROM plans WHERE id = $1`, planID).
		Scan(&p.ID, &p.Name, &p.Description, &p.PriceCents, &p.BillingPeriod, &p.MaxSeats, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Plan{}, ErrPlanNotFound
	}
	return p, err
}

func (r *Repository) UpsertSubscription(ctx context.Context, input ActivateInput, activatedAt time.Time, currentPeriodEnd time.Time) (TenantSubscription, error) {
	var sub TenantSubscription
	err := r.pool.QueryRow(ctx, `
INSERT INTO tenant_subscriptions (tenant_id, plan_id, seats, status, activated_at, current_period_end, updated_at)
VALUES ($1, $2, $3, 'active', $4, $5, NOW())
ON CONFLICT (tenant_id) DO UPDATE SET
	plan_id = EXCLUDED.plan_id,
	seats = EXCLUDED.seats,
	status = EXCLUDED.status,
	activated_at = EXCLUDED.activated_at,
	current_period_end = EXCLUDED.current_period_end,
	updated_at = NOW()
RETURNING id, tenant_id, plan_id, status, seats, activated_at, current_period_end
`, input.TenantID, input.PlanID, input.Seats, activatedAt, currentPeriodEnd).
		Scan(&sub.ID, &sub.TenantID, &sub.PlanID, &sub.Status, &sub.Seats, &sub.ActivatedAt, &sub.CurrentPeriodEnd)
	return sub, err
}

func (r *Repository) GetSubscription(ctx context.Context, tenantID string) (TenantSubscription, error) {
	var sub TenantSubscription
	err := r.pool.QueryRow(ctx, `SELECT id, tenant_id, plan_id, status, seats, activated_at, current_period_end FROM tenant_subscriptions WHERE tenant_id = $1`, tenantID).
		Scan(&sub.ID, &sub.TenantID, &sub.PlanID, &sub.Status, &sub.Seats, &sub.ActivatedAt, &sub.CurrentPeriodEnd)
	if errors.Is(err, pgx.ErrNoRows) {
		return TenantSubscription{}, ErrSubscriptionNotFound
	}
	return sub, err
}
