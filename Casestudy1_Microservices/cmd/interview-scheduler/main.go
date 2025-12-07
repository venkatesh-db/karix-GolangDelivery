package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"casestudy1_microservices/internal/events"
	"casestudy1_microservices/internal/infra/kafka"
	"casestudy1_microservices/internal/infra/nats"
)

const (
	kafkaTopic          = "shortlisted_candidates"
	natsSubject         = "notifications.interview"
	defaultKafkaBrokers = "localhost:9092"
	natsDefaultURL      = "nats://localhost:4222"
)

type InterviewScheduler struct {
	interviews               sync.Map
	kafkaConsumer            *kafka.Consumer
	natsPublisher            *nats.Publisher
	totalCandidatesProcessed int64
	totalInterviewsScheduled int64
}

func main() {
	// Load configuration
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = defaultKafkaBrokers
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = natsDefaultURL
	}

	// Initialize Kafka consumer
	consumer := kafka.NewConsumer([]string{kafkaBrokers}, kafkaTopic, "interview-scheduler")
	defer consumer.Close()

	// Initialize NATS publisher
	natsPublisher, err := nats.NewPublisher(natsURL)
	if err != nil {
		slog.Error("Failed to connect to NATS", slog.Any("error", err))
		os.Exit(1)
	}
	defer natsPublisher.Close()

	scheduler := &InterviewScheduler{
		kafkaConsumer: consumer,
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

	// Start HTTP server
	http.HandleFunc("/interviews/", scheduler.handleGetInterview)
	server := &http.Server{
		Addr:    ":8081",
		Handler: nil,
	}
	go func() {
		slog.Info("Starting HTTP server", slog.String("address", ":8081"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Start Kafka consumer loop
	scheduler.runConsumerLoop(ctx)

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("HTTP server shutdown failed", slog.Any("error", err))
	}
}

func (s *InterviewScheduler) runConsumerLoop(ctx context.Context) {
	slog.Info("Starting Kafka consumer loop")
	if err := s.kafkaConsumer.Consume(ctx, func(key, value []byte) error {
		var candidate events.ShortlistedCandidate
		if err := json.Unmarshal(value, &candidate); err != nil {
			slog.Error("Failed to deserialize message", slog.String("key", string(key)), slog.Any("error", err))
			return nil // Skip invalid messages
		}
		s.processCandidate(ctx, candidate)
		return nil
	}); err != nil {
		slog.Error("Kafka consumer error", slog.Any("error", err))
	}
}

func (s *InterviewScheduler) processCandidate(ctx context.Context, candidate events.ShortlistedCandidate) {
	// Simulate assigning an interview slot
	slot := time.Now().Add(time.Duration(rand.Intn(7)+1) * 24 * time.Hour)
	s.interviews.Store(candidate.StudentID, slot)

	// Increment metrics
	s.totalCandidatesProcessed++
	s.totalInterviewsScheduled++

	// Publish InterviewScheduled event to NATS
	event := events.InterviewScheduled{
		StudentID:     candidate.StudentID,
		Name:          candidate.Name,
		InterviewSlot: slot.Format(time.RFC3339),
	}
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("Failed to serialize InterviewScheduled event", slog.String("student_id", candidate.StudentID), slog.Any("error", err))
		return
	}
	if err := s.natsPublisher.Publish(natsSubject, data); err != nil {
		slog.Error("Failed to publish InterviewScheduled event to NATS", slog.String("student_id", candidate.StudentID), slog.Any("error", err))
		return
	}

	slog.Info("Interview scheduled", slog.String("student_id", candidate.StudentID), slog.String("slot", slot.Format(time.RFC3339)))
}

func (s *InterviewScheduler) handleGetInterview(w http.ResponseWriter, r *http.Request) {
	studentID := r.URL.Path[len("/interviews/"):]
	if studentID == "" {
		http.Error(w, "student_id is required", http.StatusBadRequest)
		return
	}

	value, ok := s.interviews.Load(studentID)
	if !ok {
		http.Error(w, "Interview not found", http.StatusNotFound)
		return
	}

	slot, ok := value.(time.Time)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"student_id": studentID,
		"slot":       slot.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
