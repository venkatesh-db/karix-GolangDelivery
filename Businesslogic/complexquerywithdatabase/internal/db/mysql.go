package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitMySQL() error {

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/smiles?parseTime=true"

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open error: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	var pingErr error
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		pingErr = sqlDB.PingContext(ctx)
		cancel()

		if pingErr == nil {
			DB = sqlDB
			log.Println("Connected to MySQL")
			return nil
		}

		log.Printf("MySQL not ready. Retrying... (%v)", pingErr)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("failed to connect to MySQL: %v", pingErr)
}

func GetDB() *sql.DB {
	return DB
}

func Createtables() error {

	usertable := `CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	ordertable := `CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);`

	_, err := DB.Exec(usertable)

	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	_, err = DB.Exec(ordertable)

	if err != nil {
		return fmt.Errorf("error creating orders table: %v", err)
	}
	return nil

}
