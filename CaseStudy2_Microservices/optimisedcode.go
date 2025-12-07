package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

/*


export KAFKA_BROKERS="localhost:9092"
go run optimisedcode.go -mode=producer

export KAFKA_BROKERS="localhost:9092"
go run optimisedcode.go -mode=consumer


docker exec -it a1053cdbc5b3  kafka-topics --bootstrap-server localhost:9092 --create --topic bench_highspeed --partitions 12 --replication-factor 1

To see output

docker exec -it a1053cdbc5b3 kafka-console-consumer --bootstrap-server localhost:9092 --topic bench_highspeed --from-beginning



*/

// --------------------------------------------------------------------

type Config struct {
	Brokers   []string
	Topic     string
	MsgCount  int
	Mode      string
	BatchSize int
}

func loadConfig() Config {
	cfg := Config{}

	brokerEnv := os.Getenv("KAFKA_BROKERS")
	if brokerEnv == "" {
		brokerEnv = "localhost:9092"
	}

	flag.StringVar(&cfg.Topic, "topic", getEnv("KAFKA_TOPIC", "bench_highspeed"), "Kafka topic")
	flag.StringVar(&cfg.Mode, "mode", getEnv("MODE", ""), "producer / consumer")
	flag.IntVar(&cfg.MsgCount, "msg-count", getEnvInt("MSG_COUNT", 1_000_000), "number of messages")
	flag.IntVar(&cfg.BatchSize, "batch", getEnvInt("BATCH_SIZE", 1000), "producer batch size")

	flag.Parse()

	cfg.Brokers = strings.Split(brokerEnv, ",")

	if cfg.Mode == "" {
		log.Fatal("MODE is required: producer or consumer")
	}

	return cfg
}

// --------------------------------------------------------------------

func main() {
	cfg := loadConfig()

	ctx, cancel := context.WithCancel(context.Background())

	// CTRL+C shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Shutting down...")
		cancel()
	}()

	switch cfg.Mode {
	case "producer":
		runProducer(ctx, cfg)
	case "consumer":
		runConsumer(ctx, cfg)
	default:
		log.Fatalf("Unknown MODE: %s", cfg.Mode)
	}
}

// --------------------------------------------------------------------
// HIGH-THROUGHPUT PRODUCER
// --------------------------------------------------------------------

func runProducer(ctx context.Context, cfg Config) {
	log.Println("Starting optimized producer")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.Brokers,
		Topic:        cfg.Topic,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: 5 * time.Millisecond,
		Async:        true,
		RequiredAcks: kafka.RequireOne,
		Compression:  kafka.Snappy,
	})
	defer writer.Close()

	// Pre-create a small payload to avoid GC overhead
	payload := []byte("resume-payload-1234567890")

	sent := 0
	start := time.Now()

	for sent < cfg.MsgCount {
		if ctx.Err() != nil {
			break
		}

		n := cfg.BatchSize
		remaining := cfg.MsgCount - sent
		if remaining < n {
			n = remaining
		}

		batch := make([]kafka.Message, n)
		for i := 0; i < n; i++ {
			batch[i] = kafka.Message{Value: payload}
		}

		if err := writer.WriteMessages(ctx, batch...); err != nil {
			log.Printf("Producer error: %v", err)
			break
		}

		sent += n
	}

	duration := time.Since(start)
	qps := float64(sent) / duration.Seconds()

	log.Printf("Produced %d messages in %v (%.2f msg/s)",
		sent, duration, qps)
}

// --------------------------------------------------------------------
// HIGH-THROUGHPUT CONSUMER
// --------------------------------------------------------------------

func runConsumer(ctx context.Context, cfg Config) {
	log.Println("Starting optimized consumer")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  "bench-consumer-group",
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB fetch
	})
	defer reader.Close()

	var consumed int64
	start := time.Now()

	for {
		if ctx.Err() != nil {
			break
		}

		_, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Consumer error: %v", err)
			break
		}

		n := atomic.AddInt64(&consumed, 1)
		if n >= int64(cfg.MsgCount) {
			break
		}
	}

	duration := time.Since(start)
	qps := float64(consumed) / duration.Seconds()

	log.Printf("Consumed %d messages in %v (%.2f msg/s)",
		consumed, duration, qps)
}

// --------------------------------------------------------------------
// UTILS
// --------------------------------------------------------------------

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	x, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return x
}
