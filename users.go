package main

import (
	"encoding/json"
	"net/http"
	"log"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	string	`json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %w", err)
		w.WriteHeader(400)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %w", err)
		w.WriteHeader(500)
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

