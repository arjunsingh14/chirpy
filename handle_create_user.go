package main

import (
	"encoding/json"
	"github.com/arjunsingh14/chirpy/internal/auth"
	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type createUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	if reqParams.Password == "" {
		respondWithError(w, 400, "Password required")
		return
	}
	hashedPassword, err := auth.HashPassword(reqParams.Password)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: reqParams.Email, HashedPassword: hashedPassword})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJson(w, 201, buildUser(user))
}

func buildUser(u database.User) user {
	return user{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}
