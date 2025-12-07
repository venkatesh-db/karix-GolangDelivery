# Microservices Case Study

This project simulates a distributed system for a hiring pipeline using microservices. It includes the following services:

1. **ResumeCollector**: Collects resumes and publishes `ResumeUploaded` events to Kafka.
2. **ShortlistingService**: Consumes `ResumeUploaded` events, applies shortlisting rules, and publishes `ShortlistedCandidate` events.
3. **InterviewScheduler**: Consumes `ShortlistedCandidate` events, schedules interviews, and publishes `InterviewScheduled` notifications to NATS.
4. **LoadGenerator**: Simulates 50,000 resume submissions to test the system.

## Prerequisites

- **Kafka** and **Zookeeper**: Ensure Kafka and Zookeeper are running.
- **NATS**: Ensure NATS server is running.
- **Go**: Install Go (1.18 or later).
- **Docker**: Install Docker for running Kafka, Zookeeper, and NATS.

## Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd Casestudy1_Microservices
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Start infrastructure services:
   ```bash
   docker-compose up -d zookeeper kafka nats
   ```

4. Create Kafka topics:
   ```bash
   kafka-topics.sh --create --topic resume_uploaded --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   kafka-topics.sh --create --topic shortlisted_candidates --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
   ```

## Running the Services

Run the services in the following order:

1. **ResumeCollector**:
   ```bash
   go run ./cmd/resume-collector
   ```

2. **ShortlistingService**:
   ```bash
   go run ./cmd/shortlisting-service
   ```

3. **InterviewScheduler**:
   ```bash
   go run ./cmd/interview-scheduler
   ```

4. **LoadGenerator** (optional):
   ```bash
   go run ./cmd/load-generator
   ```

## Environment Variables

Set the following environment variables before running the services:

- `KAFKA_BROKERS`: Kafka broker address (e.g., `localhost:9092`)
- `NATS_URL`: NATS server URL (e.g., `nats://localhost:4222`)

Example:
```bash
export KAFKA_BROKERS="localhost:9092"
export NATS_URL="nats://localhost:4222"
```

## Project Structure

```
Casestudy1_Microservices/
├── cmd/
│   ├── resume-collector/       # ResumeCollector service
│   ├── shortlisting-service/   # ShortlistingService
│   ├── interview-scheduler/    # InterviewScheduler
│   ├── load-generator/         # LoadGenerator
├── internal/
│   ├── events/                 # Event definitions
│   ├── infra/
│       ├── kafka/              # Kafka wrapper
│       ├── nats/               # NATS wrapper
├── docker-compose.yaml         # Docker Compose file for infrastructure
├── go.mod                      # Go module file
```

## Testing the System

1. Start all services.
2. Use the `LoadGenerator` to simulate resume submissions.
3. Verify logs and outputs for successful processing of events.

## Troubleshooting

- **Kafka Topic Errors**: Ensure topics are created and Kafka is running.
- **Environment Variables**: Verify `KAFKA_BROKERS` and `NATS_URL` are set correctly.
- **Service Logs**: Check logs for detailed error messages.

## License

This project is licensed under the MIT License.