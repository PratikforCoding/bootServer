package main

import (

	"net/http"
	"sort"
)
func (cfg *apiConfig)handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.getChirps()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could't retrieve chirps!")
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	responseWithJson(w, http.StatusOK, chirps)
}