package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerResetFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environments"))
		return
	}

	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed to write the database: %w", err)))
		return
	}
	cfg.fileserverHits.Store(int32(0))
	resp := fmt.Sprintf("Hits reset to: %v", cfg.fileserverHits.Load())
	w.Write([]byte(resp))
}
