package main

import _ "github.com/lib/pq"
import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/beerded/Chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits 	atomic.Int32
	db				*database.Queries
	platform		string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening postgres database: %w", err)
	}
	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform: 		platform,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetFileserverHits)

	server := &http.Server{
		Addr:		":"+port,
		Handler:	mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
