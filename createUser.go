package main

import (
	"encoding/json"
	"net/http"
	"errors"
	"github.com/PratikforCoding/chirpy.git/internal/auth"
	"github.com/PratikforCoding/chirpy.git/internal/database"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func(cfg *apiConfig)handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type returnUser struct {
		ID int `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameter")
		return
	}

	hashedPassword, err := auth.HashedPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			responseWithError(w, http.StatusConflict, "User already exists")
			return
		}

		responseWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	
	responseWithJson(w, http.StatusCreated, returnUser {
		Email: user.Email,
		ID: user.ID,
	})
}