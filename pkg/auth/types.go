package auth

import (
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/gargath/mongster/pkg/backend"
)

type Auth struct {
	b            *backend.Backend
	oAuthConfig  *oauth2.Config
	sessionStore *sessions.CookieStore
	sessionName  string
}

type Config struct {
	SessionName string
}

type Userinfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	Verified   bool   `json:"email_verified"`
}
