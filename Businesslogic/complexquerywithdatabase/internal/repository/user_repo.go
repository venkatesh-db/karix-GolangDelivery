package repository

import (
	"context"
	"errors"
	"database/sql"
	"fmt"
	"time"
	"complex-sql/internal/models"

)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(database *sql.DB) *UserRepo {
	return &UserRepo{
		db: database,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) (int64, error) {

// sql string --> harddisk 

	query := `INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Execute the query - ram data is insertd in to harddisk
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, time.Now())
	
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil

}


func (r *UserRepo) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {

	query := `SELECT id, name, email, created_at FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User

	// harddsik data is copied to ram 

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)


		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}

		if err != nil {
			return nil, err
		}


	return &user, nil
}