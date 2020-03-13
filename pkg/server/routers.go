package server

import (
	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/auth"
)

func buildRouter() *mux.Router {
	router := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.Use(auth.TokenVerifierMiddleware)

	router.HandleFunc("/", index)
	api.HandleFunc("/", apiIndex)

	return router
}
