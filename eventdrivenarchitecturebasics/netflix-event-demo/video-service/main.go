
package main

import (
	"log"
  "github.com/nats-io/nats.go"
)

func main() {

	// Connect to NATS server

	nc, err := nats.Connect(nats.DefaultURL)
	
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to the "video.uploaded" subject

	js,_:=nc.JetStream()

    js.AddStream(&nats.StreamConfig{
		Name: "VIDEO_EVENTS",
		Subjects: []string{"video.*"},
	})

   event := `{"userId":"12345","videoId":"67890","timestamp":"2024-06-01T12:00:00Z"}`
  js.Publish("video.started", []byte(event))

  log.Println("Published video.started event")

}

