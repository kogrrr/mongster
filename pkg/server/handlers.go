package server

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><title>Mongster</title> <body> <p>Hello World from Mongster</p><p><a href='/auth/login'><button>Login with Google!</button> </a> </p> </body></html>")
}
