package main

import (
	//"fmt"
	//	"html"
	"log"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr: ":8080",
	}

	log.Fatal(s.ListenAndServe())
}
