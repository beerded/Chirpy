package main

import (
	"log"
	"encoding/json"
	"net/http"
	"time"

	"github.com/beerded/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body		string		`json:"body"`
		UserID		uuid.UUID	`json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	cleanLanguage, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:		cleanLanguage,
		UserID:		params.UserID,
	})
	if err != nil {
		log.Printf("Unable to create chirp: %w", err)
		respondWithError(w, 500, "Unable to save chirp")
		return
	}
	type jsonChirp struct {
		ID			uuid.UUID 	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Body		string		`json:"body"`
		UserID		uuid.UUID	`json:"user_id"`
	}
	respondWithJSON(w, http.StatusCreated, jsonChirp{
		ID:			chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body:		chirp.Body,
		UserID:		chirp.UserID,
	})
}
