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
	jwtSecret		string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening postgres database: %v", err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform: 		platform,
		jwtSecret:		jwtSecret,
	}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirpByID)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerTokenRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerTokenRevoke)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileserverHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetFileserverHits)

	server := &http.Server{
		Addr:		":"+port,
		Handler:	mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
