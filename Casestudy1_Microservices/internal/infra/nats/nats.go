package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

// Publisher wraps NATS connection for publishing messages.
type Publisher struct {
	nc *nats.Conn
}

// NewPublisher creates a new NATS publisher.
func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Publisher{nc: nc}, nil
}

// Publish sends a message to a NATS subject.
func (p *Publisher) Publish(subject string, data []byte) error {
	return p.nc.Publish(subject, data)
}

// Close shuts down the publisher.
func (p *Publisher) Close() {
	p.nc.Close()
}

// Subscriber wraps NATS subscription for consuming messages.
type Subscriber struct {
	sub *nats.Subscription
}

// NewSubscriber creates a new NATS subscriber.
func NewSubscriber(url, subject string, handler func(msg *nats.Msg)) (*Subscriber, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	sub, err := nc.Subscribe(subject, handler)
	if err != nil {
		return nil, err
	}
	return &Subscriber{sub: sub}, nil
}

// Close unsubscribes from the subject.
func (s *Subscriber) Close() {
	if err := s.sub.Unsubscribe(); err != nil {
		log.Printf("Error unsubscribing: %v", err)
	}
}
