package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserParams struct {
	Email string`json:"email"`
}

type user struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	reqParams := createUserParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), reqParams.Email)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	} 
	respondWithJson(w, 201, buildUser(user))
}

func buildUser(u database.User) user {
	return user{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,

	}
}