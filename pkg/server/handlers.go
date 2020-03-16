package server

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><title>Mongoose</title> <body> <p>Hello World from Mongoose</p><p><a href='/auth/login'><button>Login with Google!</button> </a> </p> </body></html>")
}
