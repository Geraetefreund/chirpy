package main

import (
	"errors"
	"strings"
)

func validateAndClean(input string) (string, error) {
	const maxChirpLength = 140

	if len(input) > 140 {
		return input, errors.New("Chirp exeeds 140 chars")
	}

	profanityMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Fields(input)
	for i, w := range words {
		if _, ok := profanityMap[strings.ToLower(w)]; ok {
			words[i] = "****"
		}
	}
	cleanedInput := strings.Join(words, " ")

	return cleanedInput, nil
}
