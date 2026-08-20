package main

import (
	"net/http"
	"github.com/arjunsingh14/chirpy/internal/database"
)


func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJson(w, 200, buildChirps(chirps))
	
}

func buildChirps(c []database.Chirp) []chirp {
	chirps := []chirp{}
	for _, chirp := range c {
		chirps = append(chirps, buildChirp(chirp))
	}
	return chirps
}
 