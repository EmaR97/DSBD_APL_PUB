package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDBService handles MongoDB connections.
type MongoDBService struct {
	client *mongo.Client
}

const (
	// DefaultTimeout is the default timeout duration for MongoDB connections.
	DefaultTimeout = 10 * time.Second
)

// NewMongoDBService creates a new instance of MongoDBService.
func NewMongoDBService(databaseURL string) (*MongoDBService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURL).SetMaxPoolSize(10))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &MongoDBService{client: client}, nil
}

// GetDatabase returns a MongoDB database instance.
func (s *MongoDBService) GetDatabase(databaseName string) *mongo.Database {
	return s.client.Database(databaseName)
}

// Close closes the MongoDB connection.
func (s *MongoDBService) Close() error {
	return s.client.Disconnect(context.Background())
}
