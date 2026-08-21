package main

import (
	"encoding/json"
	"github.com/arjunsingh14/chirpy/internal/auth"
	"github.com/arjunsingh14/chirpy/internal/database"
	"net/http"
	"time"
)

type loginParams struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
}

const accessExpiresIn = time.Hour
const refreshExpiresIn = 60 * 24 * time.Hour

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

	token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, accessExpiresIn)
	if err != nil {
		respondWithError(w, 500, "Couldn't create token")
		return
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{UserID: user.ID, Token: token, ExpiresAt: time.Now().Add(refreshExpiresIn)})

	if err != nil {
		respondWithError(w, 500, "Couldn't create refresh token")
		return
	}

	respondWithJson(w, 200, loginResponse{
		userResponse: newUserResponse(user),
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

