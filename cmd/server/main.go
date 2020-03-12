package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gargath/mongopoc/pkg/server"
)

func main() {
	log.Printf("Mongo PoC version %s\n", VERSION)
	s := &server.Server{
		GracefulShutdownPeriod: time.Second * 60,
	}
	s.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	s.Shutdown()
	os.Exit(0)
}
