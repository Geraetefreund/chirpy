package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	/*
		we could use a nested loop (simple, easy to read) but O(n) or a map for O(1)
		normally, if profanities is data or too large, we create the map from a slice
		profanities := []string{"kerfuffle", "sharbert", "fornax"}
		profanityMap := make(map[string]struct{})
		for _, p := range profanities {
			profanityMap[p] = struct{}{}
		}

		below, I simply created the map inline since it was only three words.
	*/

	profanityMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Fields(params.Body)
	for i, w := range words {
		if _, ok := profanityMap[strings.ToLower(w)]; ok {
			words[i] = "****"
		}
	}
	cleanedBody := strings.Join(words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedBody,
	})
}
