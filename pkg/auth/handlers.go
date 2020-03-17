package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gargath/mongoose/pkg/entities"
	"golang.org/x/oauth2"
)

func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := a.sessionStore.Get(r, a.sessionName)
	if err != nil {
		log.Printf("WARN: Error getting session: %v", err)
	}
	state, err := RandomHexBytes(20)
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

func (a *Auth) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	s, t, err := a.decodeCallback(r)
	if err != nil {
		http.Error(w, "malformed OAuth2 callback", http.StatusBadRequest)
		log.Printf("Error decoding callback: %v", err)
		return
	}
	session, err := a.sessionStore.Get(r, a.sessionName)

	sessionState, err := a.getStateFromSession(session)
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

	user, err := a.b.FindUserBySub(userinfo.Sub)
	if err != nil {
		http.Error(w, "backend error", http.StatusInternalServerError)
		log.Printf("Error getting user with sub %s from backend: %v", userinfo.Sub, err)
		return
	}
	if user == nil {
		user = &entities.User{
			Sub:        userinfo.Sub,
			Name:       userinfo.Name,
			FamilyName: userinfo.FamilyName,
			GivenName:  userinfo.GivenName,
			IconURL:    userinfo.Picture,
			Token: &entities.Token{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				Expiry:       token.Expiry,
				TokenType:    token.TokenType,
			},
		}
		_, err := a.b.InsertUser(user)
		if err != nil {
			http.Error(w, "backend error", http.StatusInternalServerError)
			log.Printf("Error persisting user: %v", err)
			return
		}
	} else {
		t := &entities.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			Expiry:       token.Expiry,
			TokenType:    token.TokenType,
		}
		err := a.b.UpdateUserToken(userinfo.Sub, t)
		if err != nil {
			log.Printf("error updating user token: %v", err)
		}
	}

	session.Values["email"] = userinfo.Email
	session.Values["sub"] = userinfo.Sub
	session.Save(r, w)
	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
}

func (auth *Auth) SelfHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.sessionStore.Get(r, auth.sessionName)
	sub, ok := session.Values["sub"].(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}
	u, err := auth.b.FindUserBySub(sub)
	if err != nil {
		log.Printf("error fetching user info: %v", err)
		http.Error(w, "backend error", http.StatusInternalServerError)
	}
	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{}"))
		return
	}

	usersJson, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJson)
}
