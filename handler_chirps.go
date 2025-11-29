package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Geraetefreund/chirpy/internal/auth"
	"github.com/Geraetefreund/chirpy/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Body string `json:"body"`
}
type Chirp struct {
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

	// parse brearer from headers
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	// validate JWT
	userId, err := auth.ValidateJWT(tokenStr, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	// set dbParams.UserID from the token's subject/claim

	// validate + sanitize -> cleanedBody
	cleanedBody, err := validateAndClean(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	dbParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create chirp", err)
		return
	}

	response := Chirp{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "database error", err)
		return
	}

	resp := Chirp{
		ID:        chirp.ID.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}
	respondWithJSON(w, http.StatusOK, resp)

}

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
	}

	userId, err := auth.ValidateJWT(tokenStr, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
	}

	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	// check whether user is the author of the chirp
	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "database error", err)
		return
	}

	if userId != chirp.UserID {
		respondWithJSON(w, http.StatusForbidden, "fuckoff")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, "chirp deleted successfully")

}

func (cfg *apiConfig) handlerGetAllChirpsByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	var dbChirps []database.Chirp
	var err error
	id, _ := uuid.Parse(userID)

	if userID == "" {
		if sort == "desc" {
			dbChirps, err = cfg.db.GetAllChirpsDesc(r.Context())
		} else {
			dbChirps, err = cfg.db.GetChirps(r.Context())
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps from database", err)
			return
		}
	} else {
		if sort == "desc" {
			dbChirps, err = cfg.db.GetChirpsByIDDesc(r.Context(), id)
		} else {
			dbChirps, err = cfg.db.GetChirpsByID(r.Context(), id)
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps from database", err)
			return
		}
	}

	out := make([]Chirp, 0, len(dbChirps))
	for _, c := range dbChirps {
		out = append(out, Chirp{
			ID:        c.ID.String(),
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID.String(),
		})
	}
	respondWithJSON(w, http.StatusOK, out)

}
