package mongo

import (
	"context"
	"fmt"

	"github.com/venkatesh/mongodb-simulator/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect establishes the Mongo client with sane defaults.
func Connect(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	clientOpts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetRetryWrites(true).
		SetMaxPoolSize(uint64(cfg.Workers * 4)).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(cfg.ConnectTimeout)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	return client, nil
}

// EnsureIndexes adds the minimal indexes needed for production-style workloads.
func EnsureIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "txn_id", Value: 1}}, Options: options.Index().SetUnique(true).SetName("uid_txn")},
		{Keys: bson.D{{Key: "utr", Value: 1}}, Options: options.Index().SetName("idx_utr")},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}, Options: options.Index().SetName("idx_status_created")},
		{Keys: bson.D{{Key: "payer.customer_id", Value: 1}}, Options: options.Index().SetName("idx_payer")},
		{Keys: bson.D{{Key: "compliance_flags.aml_hit", Value: 1}}, Options: options.Index().SetName("idx_aml")},
	}

	if _, err := coll.Indexes().CreateMany(ctx, models); err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}

	return nil
}

// Disconnect flushes connections politely.
func Disconnect(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		return nil
	}
	return client.Disconnect(ctx)
}
