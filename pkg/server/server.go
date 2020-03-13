package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/gargath/mongoose/pkg/auth"
)

func (s *Server) Run() error {
	r := mux.NewRouter()
	sub := r.PathPrefix("/api").Subrouter()
	sub.Use(auth.TokenVerifierMiddleware)

	r.HandleFunc("/", index)
	sub.HandleFunc("/", apiIndex)

	s.srv = &http.Server{
		Addr:         s.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Printf("Server listening on %s", s.Addr)

	if err := s.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %v", err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.GracefulShutdownPeriod)
	defer cancel()
	log.Println("Beginning Shutdown...")
	s.srv.SetKeepAlivesEnabled(false)
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shut down server: %v", err)
	}
	log.Println("Shutdown complete")
	return nil
}