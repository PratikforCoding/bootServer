package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PratikforCoding/chirpy.git/internal/auth"
)

func (cfg *apiConfig)handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	type response struct {
		User
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "couldn't validate JWT")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode parameter")
		return
	}
	hashedPassword, err := auth.HashedPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	user, err := cfg.DB.UpdateUser(userIDInt, params.Email, hashedPassword)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	responseWithJson(w, http.StatusOK, response {
		User: User {
			ID: user.ID,
			Email: user.Email,
		},
	})

}