package main

import (
	"encoding/json"
	"github.com/Geraetefreund/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUpdateEmailAndPW(w http.ResponseWriter, r *http.Request) {
	// what should be inside the body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//what we return back
	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

}
func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
		return
	}
	// mark revoked_at and updated_at in DB, then 204
	rows, err := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil || rows == 0 {
		// treat missing/already revoked as unauthorized
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token: ", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}
