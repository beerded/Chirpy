package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/beerded/Chirpy/internal/auth"
	"github.com/beerded/Chirpy/internal/database"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 		string	`json:"email"`
		Password	string	`json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Cannot create user. Password not provided")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Fatalf("Could not hash password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:			params.Email,
		HashedPassword:	hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	user := User{
		ID:			dbUser.ID,
		CreatedAt:	dbUser.CreatedAt,
		UpdatedAt:	dbUser.UpdatedAt,
		Email:		dbUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		NewPassword		string	`json:"password"`
		NewEmail		string	`json:"email"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Could not find JWT: %v",err))
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid auth token: %v", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request parameters")
		return
	}

	if params.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, "Cannot update user. Password not provided")
		return
	}

	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		log.Fatalf("Could not hash password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	dbUser, err := cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		ID:				userID,
		Email:			params.NewEmail,
		HashedPassword:	hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:			dbUser.ID,
		CreatedAt:	dbUser.CreatedAt,
		UpdatedAt:	dbUser.UpdatedAt,
		Email:		dbUser.Email,
	})
}
