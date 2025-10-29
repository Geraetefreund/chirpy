package main

import (
	"encoding/json"
	"github.com/Geraetefreund/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type parameters struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type chirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	// validate + sanitize -> cleanedBody
	cleanedBody, err := validateAndClean(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}
	// of course, that was obvious NOT!
	uID, _ := uuid.Parse(params.UserID)

	dbParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: uID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create chirp", err)
		return
	}

	response := chirpResponse{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
	}
	respondWithJSON(w, http.StatusCreated, response)
}
