package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides persistence for users.
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

var ErrNotFound = errors.New("user not found")

func (r *Repository) ListByTenant(ctx context.Context, tenantID string) ([]User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, tenant_id, email, full_name, created_at FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.FullName, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

func (r *Repository) Create(ctx context.Context, input CreateInput) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx, `INSERT INTO users (tenant_id, email, full_name) VALUES ($1, $2, $3) RETURNING id, tenant_id, email, full_name, created_at`, input.TenantID, input.Email, input.FullName).
		Scan(&u.ID, &u.TenantID, &u.Email, &u.FullName, &u.CreatedAt)
	return u, err
}

func ensureTenant(ctx context.Context, pool *pgxpool.Pool, tenantID string) error {
	cmd, err := pool.Exec(ctx, `INSERT INTO tenants (id, name) VALUES ($1, $1) ON CONFLICT (id) DO NOTHING`, tenantID)
	if err != nil {
		return err
	}
	_ = cmd
	return nil
}

func (r *Repository) CreateWithTenant(ctx context.Context, input CreateInput) (User, error) {
	if err := ensureTenant(ctx, r.pool, input.TenantID); err != nil {
		return User{}, err
	}
	return r.Create(ctx, input)
}
