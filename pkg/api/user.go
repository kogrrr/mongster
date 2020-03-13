package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (api *API) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := api.b.ListUsers()
	if err != nil {
		log.Printf("failed to list users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usersJson, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		log.Printf("failed to marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersJson)
}
