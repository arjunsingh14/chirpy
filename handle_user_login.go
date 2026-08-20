package main

import (
	"encoding/json"
	"github.com/arjunsingh14/chirpy/internal/auth"
	"net/http"
	"time"
)

type loginParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	reqParams := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	setExpiration(&reqParams)
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

	token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, time.Duration(reqParams.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, 500, "Couldn't create token")
		return
	}

	respondWithJson(w, 200, loginResponse{
		userResponse: newUserResponse(user),
		Token:        token,
	})
}

func setExpiration(params *loginParams) {
	if params.ExpiresInSeconds > 3600 || params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 3600
	}
}
