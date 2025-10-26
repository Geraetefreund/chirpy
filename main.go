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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Metrics reset"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	metrics := apiConfig{}

	myServeMux := http.NewServeMux()
	myServeMux.Handle("/app/", metrics.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	myServeMux.Handle("/app/assets/", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	myServeMux.HandleFunc("/healthz/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	myServeMux.HandleFunc("/metrics/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// w.WriteHeader(http.StatusOK) // obsolete, automatically when writing to body
		fmt.Fprintf(w, "Hits: %d", metrics.fileserverHits.Load())
	})

	myServeMux.HandleFunc("/reset/", func(w http.ResponseWriter, req *http.Request) {
		metrics.resetMetrics(w, req)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: myServeMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
