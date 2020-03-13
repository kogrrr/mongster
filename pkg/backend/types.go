package backend

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Backend struct {
	m *mongo.Client
}

type BackendConfig struct {
	MongoURI          string        `default:"mongodb://localhost:27017"`
	ConnectionTimeout time.Duration `default:"2000000000"`
}
