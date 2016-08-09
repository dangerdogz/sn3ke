package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Snake struct {
	Id  string
	Url string
}

func main() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var s Snake
		if r.Body == nil {
			http.Error(w, "Empty body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		log.Printf("openWeatherMap: %s: %.2f", s.Id, 1.2)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(s)
	})
	http.Handle("/socket", websocket.Handler(socketHandler))
	http.ListenAndServe(":8080", nil)
}

func socketHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])

	m, err := ws.Write(msg[:n])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}
