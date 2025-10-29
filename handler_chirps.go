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

// was used for creatChirp... but I am not sure of that is the correct way???
type chirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	type ChirpResp struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    string    `json:"user_id"`
	}

	out := make([]ChirpResp, 0, len(dbChirps))
	for _, c := range dbChirps {
		out = append(out, ChirpResp{
			ID:        c.ID.String(),
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID.String(),
		})
	}

	respondWithJSON(w, http.StatusOK, out)
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
