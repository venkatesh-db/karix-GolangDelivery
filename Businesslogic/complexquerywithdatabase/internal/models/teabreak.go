package models
import "time"

type User struct{
	ID        int64    `json:"id"`
	Name 	string   `json:"name"`
	Email     string   `json:"email"`
	Active   bool     `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct{
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	Amount    float64  `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserwithOrder struct{
    User
	Orders []Order `json:"orders"`
}
