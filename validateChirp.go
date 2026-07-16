package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body string`json:"body"`
}

type respBody struct {
	Cleaned_body string `json:"cleaned_body"`
}

type respErr struct {
	Error string`json:"Error"`
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJson(w, 200, respBody{Cleaned_body: cleanChirp(params.Body)})
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


