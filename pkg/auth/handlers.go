package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gargath/mongster/pkg/entities"
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
			Roles:      []entities.Role{entities.Role("admin")}, //TODO: Don't make _everyone_ admin!
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

	ourT, err := a.generateToken(userinfo, []entities.Role{entities.AdminRole}) //TODO: Don't make _everyone_ admin!
	if err != nil {
		http.Error(w, "authentication error", http.StatusInternalServerError)
		log.Printf("Error creating token: %v", err)
		return
	}

	session.Values["token"] = ourT
	session.Save(r, w)
	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
}

func (auth *Auth) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.sessionStore.Get(r, auth.sessionName)
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
}

func (auth *Auth) SelfHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	session, _ := auth.sessionStore.Get(r, auth.sessionName)
	t, ok := session.Values["token"].(string)
	if !ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return auth.secret, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{}"))
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		re := &UserinfoResponse{}
		u, err := UserinfoFromClaims(claims)
		if err != nil {
			log.Printf("error converting claims to userinfo: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{}"))
			return
		}
		re.User = *u
		roles := []entities.Role{}
		for _, rs := range claims["roles"].([]interface{}) {
			r := RoleFromString(rs.(string))
			roles = append(roles, r)
		}
		re.Roles = roles
		re.Token = t
		usersJson, err := json.MarshalIndent(re, "", "\t")
		if err != nil {
			log.Printf("failed to marshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(usersJson)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{}"))
		return
	}

}
