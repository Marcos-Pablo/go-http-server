package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code >= 500 {
		log.Printf("Responding with 5XX error: %s\n", err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorResponse{
		Error: msg,
	})
}
