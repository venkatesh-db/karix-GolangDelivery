package users

import (
	"context"
	"testing"
)

type stubRepo struct {
	listTenant   string
	listResult   []User
	listErr      error
	createdInput CreateInput
	createUser   User
	createErr    error
}

func (s *stubRepo) ListByTenant(ctx context.Context, tenantID string) ([]User, error) {
	s.listTenant = tenantID
	return s.listResult, s.listErr
}

func (s *stubRepo) CreateWithTenant(ctx context.Context, input CreateInput) (User, error) {
	s.createdInput = input
	if s.createErr != nil {
		return User{}, s.createErr
	}
	return s.createUser, nil
}

func TestServiceListValidation(t *testing.T) {
	svc := NewService(&stubRepo{})
	if _, err := svc.List(context.Background(), ""); err != ErrInvalidTenant {
		t.Fatalf("expected ErrInvalidTenant, got %v", err)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	svc := NewService(&stubRepo{})
	cases := []struct {
		name  string
		input CreateInput
		want  error
	}{
		{"missing tenant", CreateInput{Email: "a@b.com", FullName: "Name"}, ErrInvalidTenant},
		{"missing email", CreateInput{TenantID: "t", FullName: "Name"}, ErrInvalidEmail},
		{"missing name", CreateInput{TenantID: "t", Email: "a@b.com"}, ErrInvalidName},
	}
	for _, tc := range cases {
		if _, err := svc.Create(context.Background(), tc.input); err != tc.want {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, err)
		}
	}
}

func TestServiceCreateSuccessNormalizesInput(t *testing.T) {
	repo := &stubRepo{createUser: User{ID: "u1"}}
	svc := NewService(repo)
	user, err := svc.Create(context.Background(), CreateInput{
		TenantID: " tenant-123 ",
		Email:    "ADMIN@EXAMPLE.COM",
		FullName: "  Admin User  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "u1" {
		t.Fatalf("expected returned user, got %+v", user)
	}
	if repo.createdInput.TenantID != "tenant-123" {
		t.Fatalf("tenant not trimmed: %q", repo.createdInput.TenantID)
	}
	if repo.createdInput.Email != "admin@example.com" {
		t.Fatalf("email not normalized: %q", repo.createdInput.Email)
	}
	if repo.createdInput.FullName != "Admin User" {
		t.Fatalf("name not trimmed: %q", repo.createdInput.FullName)
	}
}
