package server

import (
	"context"
	"fmt"
	"time"

	"github.com/energimind/identity-service/core/config"
	"github.com/energimind/identity-service/core/infra/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongoDB(cfg *config.MongoConfig) (*mongo.Client, error) {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cred := options.Credential{
		AuthSource:  cfg.Database,
		Username:    cfg.Username,
		Password:    cfg.Password,
		PasswordSet: false,
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Address).SetAuth(cred))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if pErr := client.Ping(ctx, nil); pErr != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", pErr)
	}

	logger.Info().Str("address", cfg.Address).Str("database", cfg.Database).Msg("Connected to MongoDB")

	return client, nil
}

func disconnectFromMongoDB(client *mongo.Client) {
	const timeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if dErr := client.Disconnect(ctx); dErr != nil {
		logger.Warn().Err(dErr).Msg("Failed to disconnect from MongoDB")

		return
	}

	logger.Info().Msg("Disconnected from MongoDB")
}
