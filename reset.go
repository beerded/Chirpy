package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerResetFileserverHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	a.fileserverHits.Store(int32(0))
	resp := fmt.Sprintf("Hits reset to: %v", a.fileserverHits.Load())
	w.Write([]byte(resp))
}
