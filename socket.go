package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-martini/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6000 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// map of all polls
var collection = make(map[int]*Poll)

type connection struct {
	ws   *websocket.Conn
	feed chan Comment
}

func serveWebsocket(w http.ResponseWriter, r *http.Request, params martini.Params) (int, string) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		return 403, "Origin not allowed"
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return 400, "Not a websocket handshake"
	} else if err != nil {
		log.Println(err)
		return 503, "internal server error"
	}
	c := &connection{feed: make(chan Comment, 256), ws: ws}
	// check if Poll exists
	var p *Poll
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 503, "Pass a number"
	}
	if poll, ok := collection[id]; ok {
		p = poll
	} else {
		p = &Poll{
			broadcast:   make(chan Comment),
			register:    make(chan *connection),
			unregister:  make(chan *connection),
			connections: make(map[*connection]bool),
		}
		go p.run()
		collection[id] = p
	}
	p.register <- c
	go c.writePump()
	c.readPump(p)
	return 200, "websocket connection closed"
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump(poll *Poll) {
	defer func() {
		poll.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		comment := Comment{}
		err := c.ws.ReadJSON(&comment)
		if err != nil {
			break
		}
		poll.broadcast <- comment
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(payload Comment) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case comment, ok := <-c.feed:
			if !ok {
				c.write(Comment{})
				return
			}
			if err := c.write(comment); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(Comment{}); err != nil {
				return
			}
		}
	}
}

type Comment struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
