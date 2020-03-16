package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func (a *auth) decodeCallback(r *http.Request) (string, string, error) {
	s, ok := r.URL.Query()["state"]
	if !ok || len(s) != 1 {
		return "", "", fmt.Errorf("OAuth2 Callback without state param")
	}
	t, ok := r.URL.Query()["code"]
	if !ok || len(t) != 1 {
		return "", "", fmt.Errorf("Oauth2 Callback without a token")
	}
	return s[0], t[0], nil
}

func (a *auth) getSession(r *http.Request) (*sessions.Session, error) {
	return a.sessionStore.Get(r, a.sessionName)
}

func (a *auth) getStateFromSession(r *http.Request) (string, error) {
	session, err := a.getSession(r)
	if err != nil {
		return "", fmt.Errorf("Error getting session: %v", err)
	}
	sessionState, ok := session.Values["oauth_state"].(string)
	if !ok {
		return "", fmt.Errorf("Session state not a string?! %v", session.Values["oauth_state"])
	}
	return sessionState, nil
}

func (a *auth) fetchUserInfo(token *oauth2.Token) (*Userinfo, error) {
	var userinfo *Userinfo
	client := a.oAuthConfig.Client(oauth2.NoContext, token)
	u, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("Error getting userinfo: %v", err)
	}

	defer u.Body.Close()
	dec := json.NewDecoder(u.Body)
	err = dec.Decode(&userinfo)
	if err != nil {
		return nil, fmt.Errorf("Error decoding userinfo response: %v", err)
	}
	return userinfo, nil
}

func randomHexBytes(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
