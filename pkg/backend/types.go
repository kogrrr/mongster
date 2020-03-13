package backend

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Backend struct {
	m *mongo.Client
}

type BackendConfig struct {
	MongoConnstr      string
	ConnectionTimeout time.Duration
}
