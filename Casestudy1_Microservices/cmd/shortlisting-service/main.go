package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"casestudy1_microservices/internal/events"
	"casestudy1_microservices/internal/infra/kafka"
	"casestudy1_microservices/internal/infra/nats"
)

const (
	defaultWorkerCount = 5
	defaultKafkaTopic  = "resume_uploaded"
	shortlistedTopic   = "shortlisted_candidates"
	natsSubject        = "notifications.shortlist"
)

type ShortlistingService struct {
	kafkaConsumer *kafka.Consumer
	kafkaProducer *kafka.Producer
	natsPublisher *nats.Publisher
	cache         sync.Map
}

// Define a local Message struct to replace kafka.Message
type Message struct {
	Key   []byte
	Value []byte
}

func main() {
	// Load configuration
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		slog.Error("KAFKA_BROKERS environment variable is required")
		os.Exit(1)
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		slog.Error("NATS_URL environment variable is required")
		os.Exit(1)
	}
	workerCount := defaultWorkerCount
	if wc := os.Getenv("WORKER_COUNT"); wc != "" {
		if n, err := strconv.Atoi(wc); err == nil {
			workerCount = n
		}
	}

	// Initialize Kafka consumer and producer
	consumer := kafka.NewConsumer([]string{kafkaBrokers}, defaultKafkaTopic, "shortlisting-service")
	defer consumer.Close()
	producer := kafka.NewProducer([]string{kafkaBrokers}, shortlistedTopic)
	defer producer.Close()

	// Initialize NATS publisher
	natsPublisher, err := nats.NewPublisher(natsURL)
	if err != nil {
		slog.Error("Failed to connect to NATS", slog.Any("error", err))
		os.Exit(1)
	}
	defer natsPublisher.Close()

	service := &ShortlistingService{
		kafkaConsumer: consumer,
		kafkaProducer: producer,
		natsPublisher: natsPublisher,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Shutting down gracefully")
		cancel()
	}()

	// Start worker pool
	service.runWorkerPool(ctx, workerCount)
}

func (s *ShortlistingService) runWorkerPool(ctx context.Context, workerCount int) {
	eventsChan := make(chan Message, workerCount*2)
	var wg sync.WaitGroup

	// Start Kafka consumer loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting Kafka consumer loop")
		if err := s.kafkaConsumer.Consume(ctx, func(key, value []byte) error {
			select {
			case eventsChan <- Message{Key: key, Value: value}:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		}); err != nil {
			slog.Error("Kafka consumer error", slog.Any("error", err))
		}
	}()

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slog.Info("Worker started", slog.Int("worker_id", workerID))
			for {
				select {
				case msg := <-eventsChan:
					s.processEvent(ctx, msg)
				case <-ctx.Done():
					slog.Info("Worker shutting down", slog.Int("worker_id", workerID))
					return
				}
			}
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()
	slog.Info("All workers have shut down")
}

func (s *ShortlistingService) processEvent(ctx context.Context, msg Message) {
	// Check idempotency
	if _, exists := s.cache.LoadOrStore(string(msg.Key), true); exists {
		slog.Info("Duplicate event detected, skipping", slog.String("key", string(msg.Key)))
		return
	}

	var resume events.ResumeUploaded
	if err := json.Unmarshal(msg.Value, &resume); err != nil {
		slog.Error("Failed to deserialize message", slog.String("key", string(msg.Key)), slog.Any("error", err))
		return
	}

	// Apply shortlisting rules
	shortlisted := resume.CGPA >= 8.0 || contains(resume.Skills, "golang") || contains(resume.Skills, "distributed systems")
	if shortlisted {
		slog.Info("Candidate shortlisted", slog.String("student_id", resume.StudentID), slog.Float64("cgpa", resume.CGPA))

		// Publish to Kafka
		shortlistedEvent := events.ShortlistedCandidate{
			StudentID: resume.StudentID,
			Name:      resume.Name,
			CGPA:      resume.CGPA,
			Branch:    resume.Branch,
		}
		data, err := json.Marshal(shortlistedEvent)
		if err != nil {
			slog.Error("Failed to serialize shortlisted event", slog.String("student_id", resume.StudentID), slog.Any("error", err))
			return
		}
		if err := s.kafkaProducer.Publish(ctx, resume.StudentID, data); err != nil {
			slog.Error("Failed to publish shortlisted event", slog.String("student_id", resume.StudentID), slog.Any("error", err))
			return
		}

		// Publish to NATS
		if err := s.natsPublisher.Publish(natsSubject, data); err != nil {
			slog.Error("Failed to publish NATS notification", slog.String("student_id", resume.StudentID), slog.Any("error", err))
		}
	} else {
		slog.Info("Candidate rejected", slog.String("student_id", resume.StudentID), slog.Float64("cgpa", resume.CGPA))
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
