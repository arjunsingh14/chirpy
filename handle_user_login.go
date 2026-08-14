package main

import (
	"encoding/json"
	"github.com/arjunsingh14/chirpy/internal/auth"
	"net/http"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	reqParams := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	user, err := cfg.db.GetUser(r.Context(), reqParams.Email)

	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	isValid, err := auth.CheckPassword(reqParams.Password, user.HashedPassword)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if !isValid {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	respondWithJson(w, 200, buildUser(user))
}
