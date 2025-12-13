package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to a NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Run publisher and subscriber
	go runPublisher(nc)
	runSubscriber(nc)
}

func runPublisher(nc *nats.Conn) {
	// Simple Publisher
	err := nc.Publish("updates", []byte("All is well"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Published message to 'updates' subject")
}

func runSubscriber(nc *nats.Conn) {
	// Simple Async Subscriber
	_, err := nc.Subscribe("updates", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	// Keep the connection alive
	select {}
}
