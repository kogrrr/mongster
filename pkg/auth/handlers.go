package auth

import (
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "mongoose-session")
	if err != nil {
		//TODO: use this pattern for other errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error getting session: %v", err)
		return
	}
	session.Values["oauth_state"] = "foo"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error saving session: %v", err)
		return
	}
	http.Redirect(w, r, "http://localhost:8080/", http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	return
}
