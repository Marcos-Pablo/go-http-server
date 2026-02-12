package main

import "net/http"

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
