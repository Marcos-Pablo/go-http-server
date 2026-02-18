package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
	"github.com/Marcos-Pablo/go-http-server/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, c.jwtKey)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	pathID := r.PathValue("chirpID")
	id, err := uuid.Parse(pathID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "You must provide a valid chirp id", err)
		return
	}

	chirp, err := c.queries.GetChirp(r.Context(), id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = c.queries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     id,
		UserID: userID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
