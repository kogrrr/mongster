package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gargath/mongoose/pkg/entities"
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
	w.Write([]byte("{\n\tusers:\n"))
	w.Write(usersJson)
	w.Write([]byte("\n}\n"))
}

func (api *API) InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"Content Type %s not supported \"\n}\n", r.Header.Get("Content-Type"))))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var u entities.User
	err := decoder.Decode(&u)
	if err != nil || u.Name == "" {
		if err == nil {
			err = fmt.Errorf("missing field 'name'")
		}
		log.Printf("error parsing %v as JSON: %v", r.Body, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\n\terror: \"malformed request document: %v\"\n}\n", err)))
		return
	}
	id, err := api.b.InsertUser(&u)
	if err != nil {
		log.Printf("failed to persist user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s/users/%s", api.prefix, id))
	w.WriteHeader(http.StatusCreated)
}
