package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}
	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't truncate table users", err)
	}

}

func (cfg *apiConfig) handlerTruncateUsersChirps(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not reset chirps", nil)
		return
	}

	err = cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not reset users", nil)
	}
	cfg.fileserverHits.Store(0)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	w.Write([]byte("Hits reset to 0"))

}
