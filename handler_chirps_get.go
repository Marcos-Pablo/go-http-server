package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (c *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := c.queries.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, chirps)
}

func (c *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	pathID := r.PathValue("chirpID")
	id, err := uuid.Parse(pathID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "You must provide a valid chirp id", err)
		return
	}

	dbChirp, err := c.queries.GetChirp(r.Context(), id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		},
	})
}
