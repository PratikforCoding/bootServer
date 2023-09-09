package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	user, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		responseWithError(w, 401, "Login failed")
		return
	}
	type returnUser struct {
		ID int `json:"id"`
		Email string `json:"email"`
	}
	responseWithJson(w, http.StatusOK, returnUser{
		ID: user.ID,
		Email: user.Email,
	})
}