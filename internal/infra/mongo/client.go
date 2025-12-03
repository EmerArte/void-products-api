package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/emerarteaga/products-api/internal/config"
	"github.com/emerarteaga/products-api/internal/infra/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client holds the MongoDB client and database
type Client struct {
	*mongo.Client
	Database *mongo.Database
}

// NewClient creates and connects to MongoDB
func NewClient(ctx context.Context, cfg *config.DatabaseConfig) (*Client, error) {
	logger.Info("connecting to MongoDB", "uri", cfg.URI, "database", cfg.Name)

	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetTimeout(time.Duration(cfg.Timeout) * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("successfully connected to MongoDB")

	return &Client{
		Client:   client,
		Database: client.Database(cfg.Name),
	}, nil
}

// Disconnect closes the MongoDB connection
func (c *Client) Disconnect(ctx context.Context) error {
	logger.Info("disconnecting from MongoDB")

	if err := c.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	logger.Info("successfully disconnected from MongoDB")
	return nil
}
