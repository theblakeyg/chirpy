package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	mux := http.NewServeMux()
	filepathRoot := "."
	port := ":8080"

	config := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	config.fileserverHits.Store(0)
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/metrics", config.metrics)
	mux.HandleFunc("/reset", config.reset)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	server.ListenAndServe()
}

func (a *apiConfig) metrics(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "text/plain; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(fmt.Sprintf("Hits: %v", a.fileserverHits.Load())))
}
