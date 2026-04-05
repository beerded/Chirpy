package main

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"
	"time"

	"github.com/beerded/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	const limit int = 140

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

	// Check if length is over limit
	if len(params.Body) > limit {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	if params.Body == "" {
		log.Printf("Input params struct:\n%+v", params)
		respondWithError(w, 400, "Missing body in request JSON, make sure to add 'body:'")
		return
	}
	if params.UserID == uuid.Nil {
		log.Printf("Input params struct:\n%+v", params)
		respondWithError(w, 400, "Missing UserID in request JSON, make sure to add 'user_id:'")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:		cleanLanguage(params.Body),
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

func cleanLanguage(original string) string {
	newWords := []string{}
	for _, word := range strings.Split(original, " ") {
		tinyWord := strings.ToLower(word)
		if (tinyWord == "kerfuffle") || (tinyWord == "sharbert") || (tinyWord == "fornax") {
			newWords = append(newWords, "****")
		} else {
			newWords = append(newWords, word)
		}
	}
	return strings.Join(newWords, " ")
}
