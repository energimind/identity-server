package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/energimind/identity-server/core/infra/mongo"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		cfg := mongo.Config{
			Address: mongoEnv.URI,
		}

		client, closer, err := mongo.Connect(ctx, cfg)
		defer closer()

		require.NotNil(t, client)
		require.NotNil(t, closer)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		cfg := mongo.Config{
			Address: "invalid-address",
		}

		client, closer, err := mongo.Connect(ctx, cfg)

		require.Nil(t, client)
		require.Nil(t, closer)
		require.Error(t, err)
	})
}
