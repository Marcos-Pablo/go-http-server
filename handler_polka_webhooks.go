package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Marcos-Pablo/go-http-server/internal/auth"
	"github.com/google/uuid"
)

type PolkaWebhookEventType string

const UserUpgraded PolkaWebhookEventType = "user.upgraded"

func (c *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var params parameters

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil || apiKey != c.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key", err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
		return
	}

	if params.Event != string(UserUpgraded) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = c.queries.UpgradeUserPlan(r.Context(), params.Data.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user plan", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
