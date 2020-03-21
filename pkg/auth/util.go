package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gargath/mongster/pkg/entities"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func (a *Auth) decodeCallback(r *http.Request) (string, string, error) {
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

func (a *Auth) getStateFromSession(session *sessions.Session) (string, error) {
	sessionState, ok := session.Values["oauth_state"].(string)
	if !ok {
		return "", fmt.Errorf("Session state not a string?! %v", session.Values["oauth_state"])
	}
	return sessionState, nil
}

func (a *Auth) fetchUserInfo(token *oauth2.Token) (*Userinfo, error) {
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

func RandomHexBytes(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *Auth) generateToken(info *Userinfo, roles []entities.Role) (string, error) {
	claims := MongsterClaims{
		Roles: roles,
		User:  info,
		StandardClaims: jwt.StandardClaims{
			Audience:  "mongster",
			Issuer:    "mongster",
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: (time.Now().Add(1 * time.Hour)).Unix(),
			Subject:   info.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(a.secret)
	return ss, err
}

func UserinfoFromClaims(claims map[string]interface{}) (*Userinfo, error) {
	ui := &Userinfo{}
	if claims["sub"] == nil {
		return ui, fmt.Errorf("claims without a sub: %v", claims)
	}
	uic := claims["userinfo"].(map[string]interface{})

	ui.Sub = uic["sub"].(string)
	ui.Name = uic["name"].(string)
	ui.GivenName = uic["given_name"].(string)
	ui.FamilyName = uic["family_name"].(string)
	ui.Picture = uic["picture"].(string)
	ui.Email = uic["email"].(string)
	ui.Verified = uic["email_verified"].(bool)
	return ui, nil
}

func RoleFromString(r string) entities.Role {
	switch r {
	case "admin":
		return entities.AdminRole
	default:
		return entities.NoneRole
	}
}
