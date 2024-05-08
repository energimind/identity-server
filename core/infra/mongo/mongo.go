package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect connects to a MongoDB instance. It returns a client instance or an error.
// It also returns a function to disconnect from the database.
func Connect(ctx context.Context, cfg Config) (*mongo.Client, func(), error) {
	client, err := connect(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	return client, func() {
		disconnect(ctx, client)
	}, nil
}

func connect(ctx context.Context, cfg Config) (*mongo.Client, error) {
	cred := options.Credential{
		AuthSource:  cfg.Database,
		Username:    cfg.Username,
		Password:    cfg.Password,
		PasswordSet: false,
	}

	opts := options.Client().ApplyURI(cfg.Address)

	if user := cfg.Username; user != "" {
		opts.SetAuth(cred)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if pErr := client.Ping(ctx, nil); pErr != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", pErr)
	}

	return client, nil
}

func disconnect(ctx context.Context, client *mongo.Client) {
	_ = client.Disconnect(ctx)
}
