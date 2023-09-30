package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"github.com/PratikforCoding/chirpy.git/internal/auth"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
	Author_id int `json:"author_id"`
}

func (cfg *apiConfig)handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Could,t decode parameters")
		return 
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, user.ID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	responseWithJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
		Author_id: chirp.Author_id,
	})
}

func validateChirp(body string) (string, error) {
	const max = 140
	if len(body) > max {
		return "", errors.New("Chirp is too long")
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