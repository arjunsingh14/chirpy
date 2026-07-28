package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/arjunsingh14/chirpy/internal/database"
	"github.com/google/uuid"
	"time"
)

type chirpParams struct {
	Body	string		`json:"body"`
	UserId	uuid.UUID	`json:"user_id"`
}

type chirp struct {
	ID			uuid.UUID `json:"id"`
	User_id		uuid.UUID `json:"user_id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string    `json:"email"`
	Body		string    `json:"body"`
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
	params.Body = cleanChirp(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body, UserID: params.UserId})

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	} 

	respondWithJson(w, 201, buildChirp(chirp))
}

func buildChirp(c database.Chirp) chirp {
	return chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.CreatedAt,
		Body: c.Body,
		User_id: c.UserID,
	}
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


