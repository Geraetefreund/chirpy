package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	myServeMux := http.NewServeMux()
	myServeMux.Handle("/", http.FileServer(http.Dir(".")))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: myServeMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
