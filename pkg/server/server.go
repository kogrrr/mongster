package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	GracefulShutdownPeriod time.Duration
	srv                    *http.Server
}

func (s *Server) Start() {
	r := mux.NewRouter()

	s.srv = &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Println("Server starting")

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.GracefulShutdownPeriod)
	defer cancel()
	log.Println("Beginning Shutdown...")
	s.srv.Shutdown(ctx)
	log.Println("Shutdown complete")
}
