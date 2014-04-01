package main

import "log"

type Poll struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan Comment

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	// question
	question *Question

	// channel to trigger graph updates
	update chan int

	// all the connections that need to be updated once a new value is received
	graphConnections map[*connection]bool

	// channel to register the connection on
	graphRegister chan *connection

	// channel to handle deregistering connection
	graphUnregister chan *connection
}

func (p *Poll) run() {
	for {
		select {
		case c := <-p.register:
			p.connections[c] = true
		case c := <-p.unregister:
			delete(p.connections, c)
			close(c.feed)
		case m := <-p.broadcast:
			for c := range p.connections {
				select {
				case c.feed <- m:
				default:
					close(c.feed)
					delete(p.connections, c)
				}
			}
		case <-p.update:
			for c := range p.graphConnections {
				err := c.ws.WriteJSON(p.question)
				if err != nil {
					log.Println(err)
					c.feed <- Comment{}
					delete(p.graphConnections, c)
				}
			}
		case c := <-p.graphRegister:
			p.graphConnections[c] = true
		case c := <-p.graphUnregister:
			delete(p.graphConnections, c)
			c.feed <- Comment{}
		}
	}
}
