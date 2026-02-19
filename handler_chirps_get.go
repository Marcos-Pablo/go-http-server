package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/Marcos-Pablo/go-http-server/internal/database"
	"github.com/google/uuid"
)

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (c *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}

	sortParam := r.URL.Query().Get("sort")
	sortDirection := asc // default

	if sortParam == string(desc) {
		sortDirection = desc
	}

	var dbChirps []database.Chirp
	if authorID != uuid.Nil {
		dbChirps, err = c.queries.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		dbChirps, err = c.queries.GetChirps(r.Context())
	}

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

	if sortDirection == desc {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
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
