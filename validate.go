package main

import (
	"encoding/json"
	"net/http"
	"log"
)


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const limit int = 140

	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid	bool	`json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(400)
		return
	}

	if len(params.Body) > limit {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	if params.Body == "" {
		log.Printf("Input params struct:\n%+v", params)
		respondWithError(w, 400, "Missing body in request JSON, make sure to add 'body:'")
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error:	msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
	return
}
