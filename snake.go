package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

var maxSnekId int = 0

type Snake struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

func NewSnake(ws *websocket.Conn, server *Server) *Snake {
	if ws == nil {
		panic("abort")
	}

	if server == nil {
		panic("nil server")
	}

	maxSnekId++
	ch := make(chan *Message, 100)
	doneCh := make(chan bool)

	return &Snake{maxSnekId, ws, server, ch, doneCh}
}

func (s *Snake) Conn() *websocket.Conn {
	return s.ws
}

func (s *Snake) Write(msg *Message) {
	select {
	case s.ch <- msg:
	default:
		s.server.Del(s)
		err := fmt.Errorf("client %d down.", s.id)
		s.server.Err(err)
	}
}

func (s *Snake) Done() {
	s.doneCh <- true
}

func (s *Snake) Listen() {
	go s.listenWrite()
	s.listenRead()
}
func (s *Snake) listenWrite() {
	log.Println("write client")
	for {
		select {
		case msg := <-s.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(s.ws, msg)

		case <-s.doneCh:
			s.server.Del(s)
			s.doneCh <- true
			return
		}
	}
}

func (s *Snake) listenRead() {
	log.Println("read client")
	for {
		select {
		case <-s.doneCh:
			s.server.Del(s)
			s.doneCh <- true
			return
		default:
			var msg Message
			err := websocket.JSON.Receive(s.ws, &msg)
			if err == io.EOF {
				s.doneCh <- true
			} else if err != nil {
				s.server.Err(err)
			} else {
				s.server.SendAll(&msg)
			}
		}
	}
}
