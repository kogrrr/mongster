package auth

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func (a *auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := a.getSession(r)
	if err != nil {
		//TODO: use this pattern for other errors
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		log.Printf("Error getting session: %v", err)
		return
	}
	state, err := randomHexBytes(20)
	if err != nil {
		http.Error(w, "Error generating oauth details", http.StatusInternalServerError)
		log.Printf("Error generating oauth state: %v", err)
		return
	}
	session.Values["oauth_state"] = state
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error saving session: %v", err)
		return
	}
	http.Redirect(w, r, a.oAuthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a *auth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	s, t, err := a.decodeCallback(r)
	if err != nil {
		http.Error(w, "malformed OAuth2 callback", http.StatusBadRequest)
		log.Printf("Error decoding callback: %v", err)
		return
	}

	sessionState, err := a.getStateFromSession(r)
	if err != nil {
		http.Error(w, "Error processing session", http.StatusInternalServerError)
		log.Printf("Error getting state from session: %v", err)
		return
	}

	if sessionState != s {
		http.Error(w, "Invalid session state", http.StatusUnauthorized)
		log.Printf("Invalid session state. Callback: %s, Session: %s", s, sessionState)
		return
	}

	token, err := a.oAuthConfig.Exchange(oauth2.NoContext, t)
	if err != nil {
		http.Error(w, "OAuth token exchange error", http.StatusBadRequest)
		log.Printf("Error during token exchange: %v", err)
		return
	}

	userinfo, err := a.fetchUserInfo(token)
	if err != nil {
		http.Error(w, "Error getting user info", http.StatusBadRequest)
		log.Printf("Error getting userinfo: %v", err)
		return
	}

	//TODO: Now look up user in mongo, create if not exist, update token otherwise

	w.WriteHeader(http.StatusOK)
	//	session.Values["email"] = userinfo.Email
	//	session.Values["picture_url"] = userinfo.Picture
	//	session.Save(r, w)
	fmt.Fprintf(w, "<html><title>Mongoose</title> <body> <p>Welcome</p><p>You are now logged in</p>\n")
	fmt.Fprintf(w, "<p> <img src=\"%s\" height=\"128\" width=\"128\"> %s </p> </body></html>", userinfo.Picture, userinfo.Email)
}
