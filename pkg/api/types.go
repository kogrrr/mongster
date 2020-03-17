package api

import (
	"github.com/gorilla/sessions"

	"github.com/gargath/mongoose/pkg/backend"
)

type API struct {
	b           *backend.Backend
	s           *sessions.CookieStore
	sessionName string
	prefix      string
}

type Config struct {
	Prefix      string
	SessionName string
}
