package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/beerded/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token	string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token Missing from Auth Header")
		return
	}

	token, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to retrieve refresh token from database")
	}

	if time.Now().UTC().After(token.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh Token Expired")
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Refresh Token revoked at %v", token.RevokedAt.Time))
		return
	}

	accessToken, err := auth.MakeJWT(token.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not issue JWT: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token:	accessToken,
	})
}

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh Token Missing from Auth Header")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not revoke token %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
