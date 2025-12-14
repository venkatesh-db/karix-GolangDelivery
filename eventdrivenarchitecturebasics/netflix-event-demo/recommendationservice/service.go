package recommendationservice

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type VideoEvent struct {
	VideoID  string `json:"video_id"`
	UserID   string `json:"user_id"`
	Category string `json:"category"`

	// camelCase fallback
	VideoIDCC string `json:"videoId"`
	UserIDCC  string `json:"userId"`
}

func InitServices() {

	store := NewStore()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	_, err = js.Subscribe("video.started", func(msg *nats.Msg) {

		log.Printf("RAW EVENT: %s", string(msg.Data))

		var event VideoEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Unmarshal error: %v", err)
			return
		}

		// Normalize IDs
		if event.UserID == "" {
			event.UserID = event.UserIDCC
		}
		if event.VideoID == "" {
			event.VideoID = event.VideoIDCC
		}

		// Enrich category
		if event.Category == "" {
			event.Category = ResolveCategory(event.VideoID)
		}

		if event.UserID == "" {
			log.Println("Invalid event: missing user id")
			return
		}

		store.Track(event.UserID, event.Category)

		top := store.GetRecommendation(event.UserID)

		log.Printf(
			"User=%s | Video=%s | Category=%s | TopRecommendation=%s",
			event.UserID,
			event.VideoID,
			event.Category,
			top,
		)
	})

	log.Println("Recommendation Service is running...")
	select {}
}
