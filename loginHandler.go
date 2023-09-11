package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/PratikforCoding/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type returnUser struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
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

	defaultExpiration := 24 * 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration{
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, time.Duration(params.ExpiresInSeconds) * time.Second)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't not create JWT")
		return
	}
	
	responseWithJson(w, http.StatusOK, returnUser{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}