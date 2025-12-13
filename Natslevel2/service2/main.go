package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Subscribe to subject "updates"
	_, err = nc.Subscribe("updates", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatalf("Error subscribing to subject: %v", err)
	}

	log.Println("Service 2 is listening for messages on subject 'updates'")

	// Keep the connection alive
	select {}
}
