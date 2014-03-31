package main

type Poll struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan Comment

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
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
		}
	}
}
