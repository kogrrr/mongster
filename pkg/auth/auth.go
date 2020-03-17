package auth

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/gargath/mongoose/pkg/backend"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//TODO: Ensure only one backend is shared between auth and api
func NewAuth(c *Config, store *sessions.CookieStore, backend *backend.Backend) (*Auth, error) {
	// Ensure crypto-safe PRNG is available
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}

	if err != nil {
		panic(fmt.Sprintf("Error generating session key. randomHex() failed with %#v", err))
	}

	a := &Auth{}

	a.sessionName = c.SessionName

	clientId := viper.GetString("clientId")
	clientSecret := viper.GetString("clientSecret")
	if clientId == "" || clientSecret == "" {
		return a, fmt.Errorf("Invalid OAuth client credentials.")
	}

	a.oAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	a.sessionStore = store

	a.b = backend
	return a, nil
}

func (a *Auth) AddRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", a.LoginHandler)
	authRouter.HandleFunc("/callback", a.CallbackHandler)
	authRouter.HandleFunc("/self", a.SelfHandler)
}
