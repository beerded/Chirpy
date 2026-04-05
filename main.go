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
	fileserverHits atomic.Int32
	dbQueries		*database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening postgres database: %w", err)
	}
	dbQueries := database.New(db)

	const port = "8080"
	apiCfg := &apiConfig{dbQueries: dbQueries}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetFileserverHits)

	server := &http.Server{
		Addr:		":"+port,
		Handler:	mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
