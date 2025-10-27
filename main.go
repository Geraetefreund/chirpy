package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

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

func validateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		Valid bool `json:"valid"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't read request")
		return
	}
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 500, "couldn't unmarshal parameters")
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJSON(w, 200, responseBody{
			Valid: true,
		})
	}
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	metrics := apiConfig{}

	myServeMux := http.NewServeMux()
	myServeMux.Handle("/app/", metrics.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	myServeMux.Handle("/app/assets", http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))
	myServeMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	myServeMux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		validateChirp(w, r)

	})

	myServeMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.WriteHeader(http.StatusOK) // obsolete, automatically when writing to body
		fmt.Fprintf(w, "<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", metrics.fileserverHits.Load())
	})

	myServeMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, req *http.Request) {
		metrics.resetMetrics(w, req)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: myServeMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
