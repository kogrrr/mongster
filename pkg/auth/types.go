package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/gargath/mongster/pkg/backend"
	"github.com/gargath/mongster/pkg/entities"
)

type Auth struct {
	b            *backend.Backend
	oAuthConfig  *oauth2.Config
	sessionStore *sessions.CookieStore
	sessionName  string
	secret       []byte
}

type Config struct {
	SessionName string
	Secret      []byte
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

type UserinfoResponse struct {
	User  Userinfo        `json:"user,inline"`
	Roles []entities.Role `json:"roles"`
	Token string          `json:"token"`
}

type MongsterClaims struct {
	Roles []entities.Role `json:"roles"`
	User  *Userinfo       `json:"userinfo"`
	jwt.StandardClaims
}
