package main

import (
	"encoding/json"
	"net/http"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
	"github.com/Marcos-Pablo/go-http-server/internal/database"
)

func (c *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	var params parameters

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

	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil || params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := c.queries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:    params.Email,
		Password: hashed,
		ID:       userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
