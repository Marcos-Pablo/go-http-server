package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
)

func (c *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenStr, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid refresh token", err)
		return
	}

	refreshToken, err := c.queries.GetRefreshToken(r.Context(), tokenStr)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to validate refresh token", err)
		return
	}

	if refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, c.jwtKey, time.Hour)

	respondWithJson(w, http.StatusOK, response{
		Token: token,
	})
}
