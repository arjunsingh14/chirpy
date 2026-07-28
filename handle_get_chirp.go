package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
		if err != nil {
		respondWithError(w, 404, err.Error())
		return
	} 
	
	respondWithJson(w, 200, buildChirp(chirp))
}