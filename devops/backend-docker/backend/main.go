package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	)`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			rows, err := db.Query("SELECT id, name, email FROM users")
			if err != nil {
				http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			users := []map[string]interface{}{}
			for rows.Next() {
				var id int
				var name, email string
				if err := rows.Scan(&id, &name, &email); err != nil {
					http.Error(w, "Failed to parse user", http.StatusInternalServerError)
					return
				}
				users = append(users, map[string]interface{}{"id": id, "name": name, "email": email})
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%v", users)
		} else if r.Method == http.MethodPost {
			var name, email string
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			name = r.FormValue("name")
			email = r.FormValue("email")

			if name == "" || email == "" {
				http.Error(w, "Name and email are required", http.StatusBadRequest)
				return
			}

			insertQuery := `INSERT INTO users (name, email) VALUES (?, ?)`
			_, err := db.Exec(insertQuery, name, email)
			if err != nil {
				http.Error(w, "Failed to insert user", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "User added successfully")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
