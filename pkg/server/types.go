package server

import (
	"net/http"
	"time"
)

type Server struct {
	GracefulShutdownPeriod time.Duration
	Addr                   string `default:"0.0.0.0:8080"`
	srv                    *http.Server
}
