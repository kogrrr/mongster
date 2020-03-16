package auth

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/gargath/mongoose/pkg/backend"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//TODO: Ensure only one backend is shared between auth and api
func newAuth() (*auth, error) {
	// Ensure crypto-safe PRNG is available
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}

	key, err := randomHexBytes(128)
	if err != nil {
		panic(fmt.Sprintf("Error generating session key. randomHex() failed with %#v", err))
	}

	a := &auth{}

	a.sessionName = "mongoose-session"

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

	a.sessionStore = sessions.NewCookieStore([]byte(key))

	config := &backend.BackendConfig{
		MongoConnstr:      viper.GetString("mongoConnstr"),
		ConnectionTimeout: 5 * time.Second,
	}
	backend, err := backend.NewBackend(config)
	if err != nil {
		return a, fmt.Errorf("failed to initialize auth backend: %v", err)
	}
	a.b = backend
	return a, nil
}

func AddRoutes(router *mux.Router) error {
	a, err := newAuth()
	if err != nil {
		return fmt.Errorf("Failed to initialize auth subsystem: %v", err)
	}
	a.addRoutes(router)
	return nil
}

func (a *auth) addRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", a.LoginHandler)
	authRouter.HandleFunc("/callback", a.CallbackHandler)
}
