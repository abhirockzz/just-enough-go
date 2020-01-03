package main

import "net/http"

func main() {
	http.Handle("/welcome", http.HandlerFunc(welcome))
	http.ListenAndServe(":8080", nil)
}

func welcome(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Welcome to Just Enough Go"))
}
