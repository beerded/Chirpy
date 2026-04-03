package main

import (
	"net/http"
	"log"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:		":"+port,
		Handler:	mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
