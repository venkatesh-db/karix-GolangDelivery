package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer wraps Kafka writer for producing messages.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
	}
}

// Publish sends a message to the Kafka topic.
func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

// Close shuts down the producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer wraps Kafka reader for consuming messages.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

// Consume reads messages from the Kafka topic.
func (c *Consumer) Consume(ctx context.Context, handler func(key, value []byte) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(m.Key, m.Value); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
