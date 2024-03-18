package testutil

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoEnvironment is a MongoDB test environment.
// It starts a single MongoDB container and provides a client to it.
type MongoEnvironment struct {
	idCounter atomic.Int64
	client    *mongo.Client
}

// Start starts the MongoDB container and returns a function to stop it.
func (m *MongoEnvironment) Start() context.CancelFunc {
	const containerStartupTimeout = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)

	mc, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "mongo:latest",
				ExposedPorts: []string{"27017/tcp"},
				WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(containerStartupTimeout),
			},
			Started: true,
		},
	)
	if err != nil {
		log.Fatal("Failed to start MongoDB container: ", err)
	}

	mongoHost, _ := mc.Host(ctx)
	mongoPort, _ := mc.MappedPort(ctx, "27017/tcp")

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", net.JoinHostPort(mongoHost, mongoPort.Port())))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	m.client = client

	return func() {
		defer cancel()

		if dErr := client.Disconnect(ctx); dErr != nil {
			log.Fatal("Failed to disconnect from MongoDB: ", dErr)
		}

		if tErr := mc.Terminate(ctx); tErr != nil {
			log.Fatal("Failed to terminate MongoDB container: ", tErr)
		}
	}
}

// NewInstance creates a new Mongo database.
// Make sure to call cancel function when done with the database.
func (m *MongoEnvironment) NewInstance() (*mongo.Database, context.CancelFunc) {
	const connectionTimeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	db := m.client.Database(fmt.Sprintf("testdb-%d", m.idCounter.Add(1)))

	closer := func() {
		if err := db.Drop(ctx); err != nil {
			log.Fatal("Failed to drop database: ", err)
		}

		cancel()
	}

	return db, closer
}
