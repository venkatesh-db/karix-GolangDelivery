package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"casestudy1_microservices/internal/events"
	"casestudy1_microservices/internal/infra/kafka"
)

// ResumeRequest represents the incoming HTTP request body.
type ResumeRequest struct {
	StudentID string   `json:"student_id"`
	Name      string   `json:"name"`
	CGPA      float64  `json:"cgpa"`
	Branch    string   `json:"branch"`
	Skills    []string `json:"skills"`
}

func main() {
	// Load Kafka configuration
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		slog.Error("KAFKA_BROKERS environment variable is required")
		os.Exit(1)
	}
	topic := "resume_uploaded"

	// Create Kafka producer
	producer := kafka.NewProducer([]string{kafkaBrokers}, topic)
	defer producer.Close()

	// HTTP server setup
	http.HandleFunc("/resumes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		var req ResumeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate input
		if err := validateResumeRequest(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Publish event to Kafka
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		event := events.ResumeUploaded{
			StudentID: req.StudentID,
			Name:      req.Name,
			CGPA:      req.CGPA,
			Branch:    req.Branch,
			Skills:    req.Skills,
		}
		data, err := json.Marshal(event)
		if err != nil {
			slog.Error("Failed to serialize event", slog.String("student_id", req.StudentID), slog.Float64("cgpa", req.CGPA), slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if err := producer.Publish(ctx, req.StudentID, data); err != nil {
			slog.Error("Failed to publish event", slog.String("student_id", req.StudentID), slog.Float64("cgpa", req.CGPA), slog.Any("error", err))
			http.Error(w, "Failed to queue resume", http.StatusInternalServerError)
			return
		}

		// TODO: Add DB storage for the resume here.

		response := map[string]string{
			"status":     "queued",
			"student_id": req.StudentID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(response)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	// Graceful shutdown
	go func() {
		slog.Info("Starting ResumeCollector service", slog.String("address", ":8080"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", slog.Any("error", err))
	}
}

func validateResumeRequest(req ResumeRequest) error {
	if req.StudentID == "" {
		return errors.New("student_id is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.CGPA <= 0 {
		return errors.New("cgpa must be greater than 0")
	}
	if len(req.Skills) == 0 {
		return errors.New("at least one skill is required")
	}
	return nil
}
