package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/venkatesh/mongodb-simulator/internal/config"
	"github.com/venkatesh/mongodb-simulator/internal/generator"
	mongowrap "github.com/venkatesh/mongodb-simulator/internal/mongo"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := mongowrap.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("mongo connect failed: %v", err)
	}
	defer func() {
		_ = mongowrap.Disconnect(context.Background(), client)
	}()

	collection := client.Database(cfg.Database).Collection(cfg.Collection)
	if err := mongowrap.EnsureIndexes(ctx, collection); err != nil {
		log.Fatalf("index creation failed: %v", err)
	}

	log.Printf("starting generator: total=%d workers=%d batch=%d", cfg.TotalRecords, cfg.Workers, cfg.BatchSize)

	gen := generator.New(cfg, collection)
	result, err := gen.Run(ctx)
	if err != nil {
		log.Fatalf("generation failed: %v", err)
	}

	log.Printf("finished inserting %d docs (failed:%d) in %s", result.Inserted, result.Failed, result.Duration)
}
