package backend

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewBackend(config *BackendConfig) (*Backend, error) {
	b := &Backend{}

	clientOptions := options.Client().ApplyURI(config.MongoConnstr)

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	log.Printf("Connecting to MongoDB on %s", config.MongoConnstr)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mongo connection ping failed: %s", err)
	}

	b.m = client

	return b, nil
}
