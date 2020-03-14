package auth

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

type OAuthCreds struct {
	clientId     string
	clientSecret string
}

var creds *OAuthCreds

//TODO: make random key random
var store = sessions.NewCookieStore([]byte("foo"))

func AddRoutes(router *mux.Router) error {
	creds = &OAuthCreds{
		clientId:     viper.GetString("clientId"),
		clientSecret: viper.GetString("clientSecret"),
	}
	if !creds.isValid() {
		return fmt.Errorf("invalid OAuth credentials")
	}

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", LoginHandler)
	authRouter.HandleFunc("/callback", CallbackHandler)

	return nil
}

func (c *OAuthCreds) isValid() bool {
	return c.clientId != "" && c.clientSecret != ""
}
