package main

import (
	"net/http"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
)

func (c *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid refresh token", err)
		return
	}

	err = c.queries.RevokeToken(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
