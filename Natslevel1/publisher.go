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

	// Simple Publisher
	err = nc.Publish("updates", []byte("Hello from Publisher!"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Published message to 'updates' subject")
}
