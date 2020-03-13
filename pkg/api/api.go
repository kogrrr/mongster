package api

import (
	"fmt"

	"github.com/gargath/mongoose/pkg/backend"
	"github.com/mcuadros/go-defaults"
)

func NewAPI() (*API, error) {
	api := &API{}
	config := &backend.BackendConfig{}
	defaults.SetDefaults(config)
	backend, err := backend.NewBackend(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize backend: %v", err)
	}
	api.b = backend
	return api, nil
}
