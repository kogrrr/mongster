package api

import (
	"fmt"
	"time"

	"github.com/gargath/mongoose/pkg/auth"
	"github.com/gargath/mongoose/pkg/backend"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func newFromConfig(c *Config) (*API, error) {
	api := &API{}
	api.prefix = c.Prefix
	config := &backend.BackendConfig{
		MongoConnstr:      viper.GetString("mongoConnstr"),
		ConnectionTimeout: 5 * time.Second,
	}
	backend, err := backend.NewBackend(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize api backend: %v", err)
	}
	api.b = backend
	return api, nil
}

func newAPI() (*API, error) {
	config := &Config{
		Prefix: "/api",
	}
	return newFromConfig(config)
}

func AddRoutes(router *mux.Router) error {
	a, err := newAPI()
	if err != nil {
		return fmt.Errorf("error while creating API: %v", err)
	}

	apiRouter := router.PathPrefix(a.prefix).Subrouter()
	apiRouter.Use(auth.TokenVerifierMiddleware)

	apiRouter.HandleFunc("/users", a.ListUsersHandler).Methods("GET")
	apiRouter.HandleFunc("/users", a.InsertUserHandler).Methods("POST")

	return nil
}
