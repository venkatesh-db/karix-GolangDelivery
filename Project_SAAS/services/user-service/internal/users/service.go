package users

import (
	"context"
	"errors"
	"strings"
)

type repository interface {
	ListByTenant(ctx context.Context, tenantID string) ([]User, error)
	CreateWithTenant(ctx context.Context, input CreateInput) (User, error)
}

// Service coordinates validation and persistence.
type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

var (
	ErrInvalidTenant = errors.New("tenant_id is required")
	ErrInvalidEmail  = errors.New("email is required")
	ErrInvalidName   = errors.New("full_name is required")
)

func (s *Service) List(ctx context.Context, tenantID string) ([]User, error) {
	if tenantID == "" {
		return nil, ErrInvalidTenant
	}
	return s.repo.ListByTenant(ctx, tenantID)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (User, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.FullName = strings.TrimSpace(input.FullName)
	if input.TenantID == "" {
		return User{}, ErrInvalidTenant
	}
	if input.Email == "" {
		return User{}, ErrInvalidEmail
	}
	if input.FullName == "" {
		return User{}, ErrInvalidName
	}
	return s.repo.CreateWithTenant(ctx, input)
}
