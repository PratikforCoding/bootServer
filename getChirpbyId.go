package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetChirpbyId(w http.ResponseWriter, r *http.Request) {
	chirpIdstr := chi.URLParam(r, "id")
	chirpId, err := strconv.Atoi(chirpIdstr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	dbChirp, err := cfg.DB.GetChirpById(chirpId)
	if err != nil {
		responseWithError(w, 404, "Chirp doesn't exist")
		return
	}
	responseWithJson(w, http.StatusOK, dbChirp)
}
