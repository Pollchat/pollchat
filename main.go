package main

import (
	"crypto/sha1"
	"fmt"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/gzip"

	"github.com/gorilla/websocket"
)

func main() {
	m:= martini.Classic()
	m.Use(gzip.All())
	m.Use(render.Renderer(render.Options{
		Directory : "templates",
		Layout : "layout",
		Extensions: []string{".tmpl", ".html"}, 
	}))
	// add an auth step here to redirect to /login if not authed
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/poll", func(r render.Render) {
		r.HTML(200, "poll", nil)
	})

	m.Get("/hash", func() string {
		h := sha1.New()
		h.Write([]byte("this is a hash"))
		hs := h.Sum(nil)
		return fmt.Sprintf("%x",hs)
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Get("/about", func(r render.Render) {
		r.HTML(200,"about",nil)
	})

	m.Get("/data", func (ws *websocket.Conn){
		fmt.Printf("Websocket connection handled")
	})

	m.Run()
}
