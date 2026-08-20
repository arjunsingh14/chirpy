package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arjunsingh14/chirpy/internal/auth"
	"github.com/arjunsingh14/chirpy/internal/database"
)

type chirpParams struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := chirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwt_secret)

	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	params.Body = cleanChirp(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: userID})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJson(w, 201, buildChirp(chirp))
}

func cleanChirp(chirp string) string {
	badWordSet := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	splitChirp := strings.Split(chirp, " ")
	for i, word := range splitChirp {
		_, ok := badWordSet[strings.ToLower(word)]
		if ok {
			splitChirp[i] = "****"
			continue
		}
	}
	return strings.Join(splitChirp, " ")
}
