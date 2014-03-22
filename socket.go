package main

import (
	"github.com/gorilla/websocket"
)

type socket struct {
	io.ReadWriter
	done	chan bool
}

func (s *Socket) Close(){
	s.done <- true
}
