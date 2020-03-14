package server

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/api"
	"github.com/gargath/mongoose/pkg/auth"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	err := api.AddRoutes(router)
	if err != nil {
		return nil, fmt.Errorf("error adding API routes: %v", err)
	}
	err = auth.AddRoutes(router)
	if err != nil {
		return nil, fmt.Errorf("error adding auth routes: %v", err)
	}

	router.HandleFunc("/", index)

	return router, nil
}
