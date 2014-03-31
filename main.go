package main

import (
	"log"
	"strconv"

	"github.com/codegangsta/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
)

func main() {
	questions[qcount] = GenerateQuestion("Who are you?", "Tom", "Laura", "Socks", "A ghost")
	m := martini.Classic()
	//m.Use(gzip.All())
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
	}))
	// add an auth step here to redirect to /login if not authed
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/poll/:pollNumber", func(r render.Render, params martini.Params) {
		id, err := strconv.Atoi(params["pollNumber"])
		if err != nil {
			log.Println(err)
		}
		if question, ok := questions[id]; ok {
			r.HTML(200, "poll", question)
		} else {
			// redirect
		}
	})

	m.Post("/poll/:pollNumber/:response", func(params martini.Params) {
		id, err := strconv.Atoi(params["pollNumber"])
		if err != nil {
			log.Println(err)
		}
		reponseID, err := strconv.Atoi(params["response"])
		if err != nil {
			log.Println(err)
		}
		questions[id].Responses[reponseID].Count++
		log.Println("Got here")
		questions[id].updated <- true
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Get("/about", func(r render.Render) {
		r.HTML(200, "about", nil)
	})

	m.Get("/data/:id", serveWebsocket)

	m.Get("/graph/:id", serveGraphWs)

	m.Run()
}
