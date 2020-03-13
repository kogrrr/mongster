package server

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/api"
)

func buildRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	err := api.AddRoutes(router)
	if err != nil {
		return nil, fmt.Errorf("error adding API routes: %v", err)
	}

	router.HandleFunc("/", index)

	return router, nil
}
