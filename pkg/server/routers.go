package server

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/api"
	"github.com/gargath/mongoose/pkg/auth"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(auth.TokenVerifierMiddleware)

	a, err := api.NewAPI()
	if err != nil {
		return nil, fmt.Errorf("error while creating API: %v", err)
	}

	router.HandleFunc("/", index)
	apiRouter.HandleFunc("/users", a.ListUsersHandler)

	return router, nil
}
