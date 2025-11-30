package repos

import (
	"context"
	"database/sql"
	"time"
	"karix.com/monolith/schemas"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, username, email string) (*schemas.User, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now()

	// FIXED: INSERT spelled correctly + check error immediately
	res, err := tx.ExecContext(
		ctx,
		"INSERT INTO users (username, email, created_at) VALUES (?, ?, ?)",
		username, email, now,
	)
	if err != nil {
		return nil, err
	}

	// FIXED: res will never be nil here
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Use same timestamp to avoid mismatch
	u := &schemas.User{
		ID:        int(id),
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return u, nil
}





func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*schemas.User, error) {

    rows:=r.db.QueryRowContext(ctx, "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", id)
    var user schemas.User
    var created string

	if err:=rows.Scan(&user.ID, &user.Username, &user.Email, &created); err != nil {
		return nil, err
	}
	
	user.CreatedAt, _ = time.Parse(time.RFC3339, created)
    return &user, nil
}