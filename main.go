package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewServer("/register")
	go server.Listen()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
