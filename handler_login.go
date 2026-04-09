package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beerded/Chirpy/internal/auth"
	"github.com/beerded/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}

	type response struct {
		User
		Token			string	`json:"token"`
		RefreshToken	string	`json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Password == "" {
		log.Println("Password not provided. Cannot login")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not issue access JWT: %v", err))
		return
	}
	refreshToken := auth.MakeRefreshToken()
	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:		refreshToken,
		UserID:		dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create refresh token: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			dbUser.ID,
			CreatedAt:	dbUser.CreatedAt,
			UpdatedAt:	dbUser.UpdatedAt,
			Email:		dbUser.Email,
		},
		Token:			accessToken,
		RefreshToken: 	refreshToken,
	})
}
