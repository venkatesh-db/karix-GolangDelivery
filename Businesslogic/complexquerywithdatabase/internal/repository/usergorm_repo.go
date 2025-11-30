package repository

import (
	"complex-sql/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepoGorm struct {
	db *sql.DB
}

func NewUserRepoGorm(database *sql.DB) *UserRepoGorm {
	return &UserRepoGorm{
		db: database,
	}
}

func (r *UserRepoGorm) CreateUserGorm(ctx context.Context, user *models.UserGorm) (int64, error) {
	query := `INSERT INTO user_gorms (name, email, active, created_at)
              VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Active,
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *UserRepoGorm) GetUserByIDGorm(ctx context.Context, userID int64) (*models.UserGorm, error) {

	query := `SELECT id, name, email, active, created_at 
              FROM user_gorms WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.UserGorm

	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
