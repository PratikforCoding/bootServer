package main

import (

	"net/http"
	"sort"
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
		responseWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}
	responseWithJson(w, http.StatusOK, Chirp {
		ID: dbChirp.ID,
		Body: dbChirp.Body,
	})
}

func (cfg *apiConfig)handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could't retrieve chirps!")
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
			Author_id: dbChirp.Author_id,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	responseWithJson(w, http.StatusOK, chirps)
}
