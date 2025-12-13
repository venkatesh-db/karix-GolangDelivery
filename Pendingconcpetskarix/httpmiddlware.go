package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ===== Generic HTTP Middleware =====

func CaromBatchMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("inside middleware handler")

		startTime := time.Now()

		log.Printf("match started")

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf("match completed in %v", duration)

	})

}

func copypastebabalji(w http.ResponseWriter, r *http.Request) {

	fmt.Println("inside copypastebabalji handler")

	w.Write([]byte("Ctrl C +Ctrl V balaji"))

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/f1carrombatch", copypastebabalji)

	handler := CaromBatchMiddleware(mux)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

}
