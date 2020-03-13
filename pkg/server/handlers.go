package server

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from Mongoose\n")
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\n\"message\": \"Hello World from Mongoose API\"\n}\n")
}
