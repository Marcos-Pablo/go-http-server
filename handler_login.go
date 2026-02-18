package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
)

func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	var params parameters

	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
		return
	}

	user, err := c.queries.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user information", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check password", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, c.jwtKey, expiresIn)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
