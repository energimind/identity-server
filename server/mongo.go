package server

import (
	"context"

	"github.com/energimind/go-kit/slog"
	"github.com/energimind/identity-server/core/config"
	driver "github.com/energimind/identity-server/core/infra/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func connectMongo(ctx context.Context, cfg config.MongoConfig, closer *closer) (*mongo.Database, error) {
	connect := func() (*mongo.Database, error) {
		client, disconnect, err := driver.Connect(ctx, mongoDriverConfig(cfg))
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to mongo")
		}

		closer.add(disconnect)

		slog.Info().Str("address", cfg.Address).Str("database", cfg.Database).Msg("Connected to mongo")

		return client.Database(cfg.Database), nil
	}

	return withDelayReporter[*mongo.Database](ctx, connect, "mongo-connect")
}

func mongoDriverConfig(cfg config.MongoConfig) driver.Config {
	return driver.Config{
		Address:  cfg.Address,
		Database: cfg.Database,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}
