package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8080"
	resumeEndpoint = "/resumes"
	numResumes     = 50000
	numWorkers     = 50
)

type Resume struct {
	StudentID string   `json:"student_id"`
	Name      string   `json:"name"`
	CGPA      float64  `json:"cgpa"`
	Branch    string   `json:"branch"`
	Skills    []string `json:"skills"`
}

func main() {
	baseURL := os.Getenv("RESUME_API_BASE")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	var totalRequests int64
	var succeeded int64
	var failed int64

	start := time.Now()

	resumes := generateResumes(numResumes)
	jobs := make(chan Resume, numResumes)
	wg := &sync.WaitGroup{}

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(baseURL, jobs, &totalRequests, &succeeded, &failed)
		}()
	}

	// Enqueue jobs
	for _, resume := range resumes {
		jobs <- resume
	}
	close(jobs)

	// Wait for workers to finish
	wg.Wait()
	duration := time.Since(start)

	// Print summary
	qps := float64(totalRequests) / duration.Seconds()
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Succeeded: %d\n", succeeded)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total Duration: %s\n", duration)
	fmt.Printf("QPS: %.2f\n", qps)
}

func generateResumes(n int) []Resume {
	branches := []string{"CSE", "ECE", "MECH", "CIVIL"}
	skills := []string{"golang", "java", "python", "distributed systems", "react", "kubernetes"}
	resumes := make([]Resume, n)

	for i := 0; i < n; i++ {
		resumes[i] = Resume{
			StudentID: fmt.Sprintf("S-%d", rand.Intn(1000000)),
			Name:      fmt.Sprintf("Student-%d", i+1),
			CGPA:      5.0 + rand.Float64()*(10.0-5.0),
			Branch:    branches[rand.Intn(len(branches))],
			Skills:    randomSkills(skills),
		}
	}

	return resumes
}

func randomSkills(skills []string) []string {
	n := rand.Intn(len(skills)) + 1
	selected := make([]string, n)
	perm := rand.Perm(len(skills))

	for i := 0; i < n; i++ {
		selected[i] = skills[perm[i]]
	}

	return selected
}

func worker(baseURL string, jobs <-chan Resume, totalRequests, succeeded, failed *int64) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for resume := range jobs {
		atomic.AddInt64(totalRequests, 1)

		data, err := json.Marshal(resume)
		if err != nil {
			slog.Error("Failed to marshal resume", slog.Any("error", err))
			atomic.AddInt64(failed, 1)
			continue
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, baseURL+resumeEndpoint, bytes.NewBuffer(data))
		if err != nil {
			slog.Error("Failed to create request", slog.Any("error", err))
			atomic.AddInt64(failed, 1)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Request failed", slog.Any("error", err))
			atomic.AddInt64(failed, 1)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			atomic.AddInt64(succeeded, 1)
		} else {
			slog.Error("Request failed with status", slog.Int("status", resp.StatusCode))
			atomic.AddInt64(failed, 1)
		}
	}
}
