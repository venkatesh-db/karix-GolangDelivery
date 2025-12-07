package main

/*


export KAFKA_BROKERS="localhost:9092"
go run main.go -mode=producer

export KAFKA_BROKERS="localhost:9092"
go run main.go -mode=consumer


docker exec -it a1053cdbc5b3  kafka-topics --bootstrap-server localhost:9092 --create --topic bench_highspeed --partitions 12 --replication-factor 1

To see output

docker exec -it a1053cdbc5b3 kafka-console-consumer --bootstrap-server localhost:9092 --topic bench_highspeed --from-beginning



*/

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	KafkaBrokers    string
	KafkaTopic      string
	MsgCount        int
	ProducerWorkers int
	ConsumerWorkers int
	Mode            string
}

func loadConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.KafkaBrokers, "brokers", os.Getenv("KAFKA_BROKERS"), "Kafka brokers (comma-separated)")
	flag.StringVar(&cfg.KafkaTopic, "topic", getEnv("KAFKA_TOPIC", "bench_highspeed"), "Kafka topic")
	flag.IntVar(&cfg.MsgCount, "msg-count", getEnvInt("MSG_COUNT", 1000000), "Number of messages")
	flag.IntVar(&cfg.ProducerWorkers, "producer-workers", getEnvInt("PRODUCER_WORKERS", 100), "Number of producer workers")
	flag.IntVar(&cfg.ConsumerWorkers, "consumer-workers", getEnvInt("CONSUMER_WORKERS", 100), "Number of consumer workers")
	flag.StringVar(&cfg.Mode, "mode", os.Getenv("MODE"), "Mode: producer, consumer, benchmark")
	flag.Parse()
	return cfg
}

func main() {
	cfg := loadConfig()

	if cfg.KafkaBrokers == "" || cfg.Mode == "" {
		log.Fatal("Kafka brokers and mode are required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	switch cfg.Mode {
	case "producer":
		runProducer(ctx, cfg)
	case "consumer":
		runConsumer(ctx, cfg)
	case "benchmark":
		runBenchmark(ctx, cfg)
	default:
		log.Fatalf("Unknown mode: %s", cfg.Mode)
	}
}

func runProducer(ctx context.Context, cfg Config) {
	log.Println("Starting producer")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      splitBrokers(cfg.KafkaBrokers),
		Topic:        cfg.KafkaTopic,
		BatchSize:    1000,
		BatchTimeout: 10 * time.Millisecond,
		Balancer:     &kafka.Murmur2Balancer{},
	})
	defer writer.Close()

	jobs := make(chan int, 10000)
	var wg sync.WaitGroup

	for i := 0; i < cfg.ProducerWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				msg := kafka.Message{
					Value: createMessage(job),
				}
				if err := writer.WriteMessages(ctx, msg); err != nil {
					log.Printf("Failed to write message: %v", err)
				}
			}
		}()
	}

	start := time.Now()
	for i := 0; i < cfg.MsgCount; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	duration := time.Since(start)
	log.Printf("Produced %d messages in %v (%.2f msg/s)", cfg.MsgCount, duration, float64(cfg.MsgCount)/duration.Seconds())
}

func runConsumer(ctx context.Context, cfg Config) {
	log.Println("Starting consumer")
	var consumed int64
	var wg sync.WaitGroup

	for i := 0; i < cfg.ConsumerWorkers; i++ {
		wg.Add(1)
		go func(partition int) {
			defer wg.Done()
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: splitBrokers(cfg.KafkaBrokers),
				Topic:   cfg.KafkaTopic,
				GroupID: "benchmark-group",
			})
			defer r.Close()

			for {
				_, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("Failed to read message: %v", err)
					return
				}
				atomic.AddInt64(&consumed, 1)
				if atomic.LoadInt64(&consumed) >= int64(cfg.MsgCount) {
					return
				}
			}
		}(i)
	}

	wg.Wait()
	log.Printf("Consumed %d messages", consumed)
}

func runBenchmark(ctx context.Context, cfg Config) {
	log.Println("Starting benchmark")
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runConsumer(ctx, cfg)
	}()

	time.Sleep(1 * time.Second)
	runProducer(ctx, cfg)

	wg.Wait()
}

func createMessage(id int) []byte {
	msg := map[string]interface{}{
		"id": id,
		"ts": time.Now().UnixNano(),
	}
	data, _ := json.Marshal(msg)
	return data
}

func splitBrokers(brokers string) []string {
	return strings.Split(brokers, ",")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
