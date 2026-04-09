package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beerded/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
		ExpiresIn	int64	`json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token		string	`json:"token"`
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

	// default timeout is 1 hour. If someone sets timeout to be longer than an
	// hour, set the duration to 1 hour instead
	expiration := 1 * time.Hour
	if params.ExpiresIn > 0 && params.ExpiresIn < 3600 {
		expiration = time.Duration(params.ExpiresIn)*time.Second
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expiration)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not issue JWT: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:			dbUser.ID,
			CreatedAt:	dbUser.CreatedAt,
			UpdatedAt:	dbUser.UpdatedAt,
			Email:		dbUser.Email,
		},
		Token:	accessToken,
	})
}
