package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/api"
	"github.com/gargath/mongoose/pkg/auth"
	"github.com/gargath/mongoose/pkg/static"
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

	router.Handle("/{rest}", http.StripPrefix("/", http.FileServer(static.Assets)))

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(static.Assets)))

	return router, nil
}
