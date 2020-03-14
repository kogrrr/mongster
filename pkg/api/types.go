package api

import (
	"github.com/gargath/mongoose/pkg/backend"
)

type API struct {
	b      *backend.Backend
	prefix string
}

type Config struct {
	Prefix string
}
