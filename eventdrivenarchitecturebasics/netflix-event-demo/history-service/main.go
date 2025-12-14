package main

import
	 (	
  "log"
  "github.com/nats-io/nats.go"

	 )
func main(){

	// Connect to NATS server
	
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to the "video.uploaded" subject

	js,_:=nc.JetStream()

	js.Subscribe("video.started", func(m *nats.Msg) {
		log.Printf("Received video.started event: %s", string(m.Data))
		// Here you would typically record the event in the user's watch history
		m.Ack()
	}, nats.Durable("HISTORY_SERVICE"), nats.ManualAck())
	
	log.Println("Subscribed to video.started events")
	
}