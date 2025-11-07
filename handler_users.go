package main

import (
	"encoding/json"
	"github.com/Geraetefreund/chirpy/internal/auth"
	"github.com/Geraetefreund/chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// what we expect from the request
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create hash", err)
		return
	}

	dbParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusConflict, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.LookUpUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	passwordOK, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "wrong password", err)
		return

	}

	if !passwordOK {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}
	var secs int
	if params.ExpiresInSeconds == nil {
		secs = 3600
	} else {
		v := *params.ExpiresInSeconds
		if v <= 0 {
			secs = 3600
		} else if v > 3600 {
			secs = 3600
		} else {
			secs = v
		}
	}

	expires := time.Duration(secs) * time.Second

	token, err := auth.MakeJWT(user.ID, cfg.secret, expires)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token: ", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})
}
