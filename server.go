package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	pattern   string
	messages  []*Message
	snakes    map[int]*Snake
	addCh     chan *Snake
	delCh     chan *Snake
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

func NewServer(pattern string) *Server {
	messages := []*Message{}
	snakes := make(map[int]*Snake)
	addCh := make(chan *Snake)
	delCh := make(chan *Snake)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		snakes,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(sc *Snake) {
	s.addCh <- sc
}

func (s *Server) Del(sc *Snake) {
	s.delCh <- sc
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.snakes {
		c.Write(msg)
	}
}

func (s *Server) Listen() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		snake := NewSnake(ws, s)
		s.Add(snake)
		snake.Listen()
	}

	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("new handler")

	for {
		select {
		case c := <-s.addCh:
			s.snakes[c.id] = c
			log.Println("Now", len(s.snakes), "snakes connected.")
		case c := <-s.delCh:
			log.Println("delete")
			delete(s.snakes, c.id)
		case msg := <-s.sendAllCh:
			log.Println("send all", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)
		case err := <-s.errCh:
			log.Println(err.Error())
		case <-s.doneCh:
			return
		}
	}
}
