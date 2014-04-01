package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/gorilla/websocket"
)

func serveGraphWs(w http.ResponseWriter, r *http.Request, params martini.Params) (int, string) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		return 403, "Origin not allowed"
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	log.Println("in graph web socket ")
	if _, ok := err.(websocket.HandshakeError); ok {
		return 400, "Not a websocket handshake"
	} else if err != nil {
		log.Println(err)
		return 503, "internal server error"
	}
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 503, "Pass a number"
	}
	// write to the socket when counts are updated
	c := &connection{ws: ws, feed: make(chan Comment, 256)}
	collection[id].graphRegister <- c
	<-c.feed
	close(c.feed)
	return 200, "websocket connection closed"
}
