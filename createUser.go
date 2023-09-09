package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Password string `json:"password"`
	Email string `json:"email"`
	ID int `json:"id"`
}

func(cfg *apiConfig)handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameter")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	type returnUser struct {
		ID int `json:"id"`
		Email string `json:"email"`
	}
	responseWithJson(w, http.StatusCreated, returnUser {
		ID: user.ID,
		Email: user.Email,
	})
}