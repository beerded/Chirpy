package main

import (
	"fmt"
	"sync/atomic"
	"net/http"
	"log"
)

func main() {
	const port = "8080"
	apiCfg := &apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerFileserverHits)
	mux.HandleFunc("/reset", apiCfg.handlerResetFileserverHits)

	server := &http.Server{
		Addr:		":"+port,
		Handler:	mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (a *apiConfig) handlerFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	resp := fmt.Sprintf("Hits: %v", a.fileserverHits.Load())
	w.Write([]byte(resp))
}

func (a *apiConfig) handlerResetFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	a.fileserverHits.Store(int32(0))
	resp := fmt.Sprintf("Hits reset to: %v", a.fileserverHits.Load())
	w.Write([]byte(resp))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
