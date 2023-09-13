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
	}

	type returnUser struct {
		ID int `json:"id"`
		Email string `json:"email"`
		AccessToken string `json:"access-token"`
		RefreshToken string `json:"refresh-token"`
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

	accessTokenTime := 60 * 60
	refreshTokenTime := 60 * 24 * 60 * 60

	accessToken, err := auth.MakeAccessToken(user.ID, cfg.jwt_secret, time.Duration(accessTokenTime) * time.Second)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't not create Access Token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken(user.ID, cfg.jwt_secret, time.Duration(refreshTokenTime) * time.Second)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't not create Access Token")
		return
	}

	responseWithJson(w, http.StatusOK, returnUser{
		
		ID: user.ID,
		Email: user.Email,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	})
}