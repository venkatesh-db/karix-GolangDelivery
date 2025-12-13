package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

// run both the servcie separtely 
// run the servcies make it effectoively
// Service 1: Publisher
// To publish a message, use the following curl command:
// curl -X POST "http://localhost:8081/publish?message=HelloNATS"



func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		message := r.URL.Query().Get("message")
		if message == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}

		// Publish message to NATS
		if err := nc.Publish("updates", []byte(message)); err != nil {
			http.Error(w, fmt.Sprintf("Error publishing message: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message published successfully"))
	})

	log.Println("Service 1 is running on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
