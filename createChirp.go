package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig)handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could,t decode parameters")
		return 
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	responseWithJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
	})
}

func validateChirp(body string) (string, error) {
	const max = 140
	if len(body) > max {
		return "", errors.New("Chirp is too long!")
	}
	badWords := map[string]struct{} {
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}