package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
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
