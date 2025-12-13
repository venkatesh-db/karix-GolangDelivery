package subscriptions

import (
	"context"
	"errors"
	"strings"
	"time"
)

type repository interface {
	ListPlans(ctx context.Context) ([]Plan, error)
	GetPlan(ctx context.Context, planID string) (Plan, error)
	UpsertSubscription(ctx context.Context, input ActivateInput, activatedAt time.Time, currentPeriodEnd time.Time) (TenantSubscription, error)
	GetSubscription(ctx context.Context, tenantID string) (TenantSubscription, error)
}

// Service validates inputs and delegates to the repository.
type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrInvalidTenantID = errors.New("tenant_id is required")
	ErrInvalidPlanID   = errors.New("plan_id is required")
	ErrInvalidSeats    = errors.New("seats must be greater than zero")
	ErrSeatLimit       = errors.New("requested seats exceed plan limit")
)

func (s *Service) Plans(ctx context.Context) ([]Plan, error) {
	return s.repo.ListPlans(ctx)
}

func (s *Service) Activate(ctx context.Context, input ActivateInput) (TenantSubscription, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.PlanID = strings.TrimSpace(input.PlanID)
	if input.TenantID == "" {
		return TenantSubscription{}, ErrInvalidTenantID
	}
	if input.PlanID == "" {
		return TenantSubscription{}, ErrInvalidPlanID
	}
	if input.Seats <= 0 {
		return TenantSubscription{}, ErrInvalidSeats
	}
	plan, err := s.repo.GetPlan(ctx, input.PlanID)
	if err != nil {
		return TenantSubscription{}, err
	}
	if plan.MaxSeats > 0 && input.Seats > plan.MaxSeats {
		return TenantSubscription{}, ErrSeatLimit
	}
	activatedAt := time.Now().UTC()
	periodEnd := activatedAt.Add(30 * 24 * time.Hour)
	return s.repo.UpsertSubscription(ctx, input, activatedAt, periodEnd)
}

func (s *Service) Subscription(ctx context.Context, tenantID string) (TenantSubscription, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return TenantSubscription{}, ErrInvalidTenantID
	}
	return s.repo.GetSubscription(ctx, tenantID)
}
