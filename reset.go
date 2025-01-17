package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) reset(responseWriter http.ResponseWriter, request *http.Request) {
	a.fileserverHits.Store(0)
	responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(fmt.Sprintf("Hits: %v", a.fileserverHits.Load())))
}
