package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
)

func main() {
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
			// should redirect here maybe 404?
			log.Println(err)
			r.Redirect("/")
		}
		if poll, ok := collection[id]; ok {
			r.HTML(200, "poll", poll.question)
		} else {
			// redirect maybe 404
			r.Redirect("/")
		}
	})

	// create a poll
	m.Post("/poll", func(req *http.Request, r render.Render, params martini.Params) {
		id := len(collection) + 1
		// parse form
		err := req.ParseForm()
		if err != nil {
			r.Redirect("/")
		}
		p := &Poll{
			broadcast:        make(chan Comment),
			register:         make(chan *connection),
			unregister:       make(chan *connection),
			connections:      make(map[*connection]bool),
			graphConnections: make(map[*connection]bool),
			graphRegister:    make(chan *connection),
			graphUnregister:  make(chan *connection),
			update:           make(chan int),
			question:         GenerateQuestion(id, req.PostForm.Get("pollquestion"), req.PostForm.Get("pollresponse1"), req.PostForm.Get("pollresponse2"), req.PostForm.Get("pollresponse3"), req.PostForm.Get("pollresponse4")),
		}
		go p.run()
		collection[id] = p
		r.Redirect("/poll/" + strconv.Itoa(id))
	})

	m.Post("/poll/:pollNumber/:response", func(params martini.Params) {
		id, err := strconv.Atoi(params["pollNumber"])
		if err != nil {
			log.Println(err)
		}
		// TODO: check the response and id
		collection[id].question.Responses[params["response"]].Count++
		// send an update to change the charts
		collection[id].update <- id
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
