package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/ivamshky/rssaggregator-go/users"
)

func HandleReadiness(rw http.ResponseWriter, r *http.Request) {
	respondJSON(rw, http.StatusOK, struct{}{})
}

func HandleError(rw http.ResponseWriter, r *http.Request) {
	respondError(rw, 400, "Something went terribly wrong")
}

func (apiCfg *apiConfig) HandleCreateUser(rw http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondError(rw, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	userRepository := users.NewPostgresSQLRepository(apiCfg.DB)
	user, err := userRepository.Create(apiCfg.ctx, users.User{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		respondError(rw, http.StatusInternalServerError, fmt.Sprintf("Error creating user: %v", err))
	}

	respondJSON(rw, http.StatusCreated, user)
}
