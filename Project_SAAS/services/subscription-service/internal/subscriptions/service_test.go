package subscriptions

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubSubscriptionRepo struct {
	plans             []Plan
	listErr           error
	planMap           map[string]Plan
	planErr           error
	lastActivateInput ActivateInput
	lastActivatedAt   time.Time
	lastPeriodEnd     time.Time
	subscription      TenantSubscription
	subscriptionErr   error
}

func (s *stubSubscriptionRepo) ListPlans(ctx context.Context) ([]Plan, error) {
	return s.plans, s.listErr
}

func (s *stubSubscriptionRepo) GetPlan(ctx context.Context, planID string) (Plan, error) {
	if s.planErr != nil {
		return Plan{}, s.planErr
	}
	if s.planMap == nil {
		return Plan{}, errors.New("plan map not set")
	}
	plan, ok := s.planMap[planID]
	if !ok {
		return Plan{}, ErrPlanNotFound
	}
	return plan, nil
}

func (s *stubSubscriptionRepo) UpsertSubscription(ctx context.Context, input ActivateInput, activatedAt time.Time, currentPeriodEnd time.Time) (TenantSubscription, error) {
	s.lastActivateInput = input
	s.lastActivatedAt = activatedAt
	s.lastPeriodEnd = currentPeriodEnd
	if s.subscriptionErr != nil {
		return TenantSubscription{}, s.subscriptionErr
	}
	if s.subscription.ID == "" {
		s.subscription = TenantSubscription{ID: "sub", TenantID: input.TenantID, PlanID: input.PlanID, Seats: input.Seats}
	}
	return s.subscription, nil
}

func (s *stubSubscriptionRepo) GetSubscription(ctx context.Context, tenantID string) (TenantSubscription, error) {
	if s.subscriptionErr != nil {
		return TenantSubscription{}, s.subscriptionErr
	}
	if s.subscription.TenantID != tenantID {
		return TenantSubscription{}, ErrSubscriptionNotFound
	}
	return s.subscription, nil
}

func TestServiceActivateValidation(t *testing.T) {
	repo := &stubSubscriptionRepo{planMap: map[string]Plan{"basic": {ID: "basic", MaxSeats: 10}}}
	svc := NewService(repo)
	cases := []struct {
		name  string
		input ActivateInput
		want  error
	}{
		{"missing tenant", ActivateInput{PlanID: "basic", Seats: 1}, ErrInvalidTenantID},
		{"missing plan", ActivateInput{TenantID: "t", Seats: 1}, ErrInvalidPlanID},
		{"invalid seats", ActivateInput{TenantID: "t", PlanID: "basic", Seats: 0}, ErrInvalidSeats},
	}
	for _, tc := range cases {
		if _, err := svc.Activate(context.Background(), tc.input); err != tc.want {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, err)
		}
	}
}

func TestServiceActivateSeatLimit(t *testing.T) {
	repo := &stubSubscriptionRepo{planMap: map[string]Plan{"basic": {ID: "basic", MaxSeats: 5}}}
	svc := NewService(repo)
	if _, err := svc.Activate(context.Background(), ActivateInput{TenantID: "t", PlanID: "basic", Seats: 6}); err != ErrSeatLimit {
		t.Fatalf("expected ErrSeatLimit, got %v", err)
	}
}

func TestServiceActivateSuccess(t *testing.T) {
	repo := &stubSubscriptionRepo{
		planMap:      map[string]Plan{"growth": {ID: "growth", MaxSeats: 100}},
		subscription: TenantSubscription{ID: "sub1", TenantID: "tenant-1", PlanID: "growth", Seats: 5},
	}
	svc := NewService(repo)
	sub, err := svc.Activate(context.Background(), ActivateInput{TenantID: " tenant-1 ", PlanID: " growth ", Seats: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.ID != "sub1" {
		t.Fatalf("expected repository subscription, got %+v", sub)
	}
	if repo.lastActivateInput.TenantID != "tenant-1" {
		t.Fatalf("tenant not trimmed: %q", repo.lastActivateInput.TenantID)
	}
	if repo.lastActivateInput.PlanID != "growth" {
		t.Fatalf("plan not trimmed: %q", repo.lastActivateInput.PlanID)
	}
	if repo.lastActivateInput.Seats != 5 {
		t.Fatalf("unexpected seats: %d", repo.lastActivateInput.Seats)
	}
	if repo.lastActivatedAt.IsZero() || repo.lastPeriodEnd.IsZero() {
		t.Fatalf("expected timestamps to be set")
	}
	if repo.lastPeriodEnd.Sub(repo.lastActivatedAt) < 30*24*time.Hour {
		t.Fatalf("period end should be at least 30 days ahead")
	}
}

func TestServiceSubscriptionValidation(t *testing.T) {
	repo := &stubSubscriptionRepo{}
	svc := NewService(repo)
	if _, err := svc.Subscription(context.Background(), ""); err != ErrInvalidTenantID {
		t.Fatalf("expected ErrInvalidTenantID, got %v", err)
	}
}
