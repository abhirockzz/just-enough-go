package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", home{})

	mux.HandleFunc("/posts", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Visit http://bit.ly/just-enough-go to get started"))
	})

	server := http.Server{Addr: ":8080", Handler: mux}

	log.Fatal(server.ListenAndServe())
}

type home struct{}

func (h home) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to the \"Just Enough Go\" blog series!!"))
}
