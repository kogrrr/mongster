package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"

	"github.com/gargath/mongoose/pkg/api"
	"github.com/gargath/mongoose/pkg/auth"
	"github.com/gargath/mongoose/pkg/backend"
	"github.com/gargath/mongoose/pkg/static"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()

	key, err := auth.RandomHexBytes(128)
	if err != nil {
		return router, fmt.Errorf("Failed to generate session key: %v", err)
	}
	s := sessions.NewCookieStore([]byte(key))

	bc := &backend.BackendConfig{
		MongoConnstr:      viper.GetString("mongoConnstr"),
		ConnectionTimeout: 5 * time.Second,
	}
	b, err := backend.NewBackend(bc)
	if err != nil {
		return router, fmt.Errorf("failed to initialize backend: %v", err)
	}

	ac := &api.Config{
		Prefix:      "/api",
		SessionName: "mongoose-session",
	}
	api, err := api.NewFromConfig(ac, s, b)
	if err != nil {
		return router, fmt.Errorf("failed to initialize api: %v", err)
	}

	auc := &auth.Config{
		SessionName: "mongoose-session",
	}
	auth, err := auth.NewAuth(auc, s, b)
	if err != nil {
		return router, fmt.Errorf("failed to initialize auth: %v", err)
	}

	api.AddRoutes(router)
	auth.AddRoutes(router)

	router.Handle("/{rest}", http.StripPrefix("/", http.FileServer(static.Assets)))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(static.Assets)))

	return router, nil
}
