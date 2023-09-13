package main

import (
	"net/http"

	"github.com/PratikforCoding/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find JWT")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't check session")
		return
	}
	if isRevoked {
		responseWithError(w, http.StatusUnauthorized, "Refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwt_secret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	responseWithJson(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Couldn't find JWT")
		return
	}

	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	responseWithJson(w, http.StatusOK, struct{}{})
}