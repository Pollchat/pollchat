package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	//"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func main() {
	m := martini.Classic()
	//m.Use(gzip.All())
	m.Use(render.Renderer(render.Options{
		Directory:       "templates",
		Layout:          "layout",
		Extensions:      []string{".tmpl", ".html"},
		HTMLContentType: "application/xhtml+xml",
	}))
	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("pollchatsecretkey"))))
	m.Use(oauth2.Google(&oauth2.Options{
		ClientId:     pollchat.clientId,
		ClientSecret: pollchat.clientSecret,
		RedirectURL:  "http://pollchat.co.uk/oauth2callback",
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
	}))
	// add an auth step here to redirect to /login if not authed
	m.Get("/", func(r render.Render, tokens oauth2.Tokens) {
		if tokens.IsExpired() {
			r.Redirect("/login")
		}
		type recent struct {
			Q   *Question
			Top string
		}
		questions := make([]*recent, 0) //questions = append(questions, question)
		if len(collection) > 5 {
			// take the top 5
			for last := len(collection); last > (len(collection) - 5); last-- {
				q := collection[last].question
				topres := ""
				currentTop := -1
				for _, res := range q.Responses {
					if res.Count > currentTop {
						currentTop = res.Count
						topres = res.Response
					}
				}
				ques := recent{collection[last].question, topres}
				questions = append(questions, &ques)
			}
		} else {
			for i := len(collection); i > 0; i-- {
				q := collection[i].question
				topres := ""
				currentTop := -1
				for _, res := range q.Responses {
					if res.Count > currentTop {
						currentTop = res.Count
						topres = res.Response
					}
				}
				ques := recent{collection[i].question, topres}
				questions = append(questions, &ques)
			}
		}
		log.Println(len(questions))
		if len(questions) != 0 {
			r.HTML(200, "indexwithrecentpolls", questions)
		} else {
			r.HTML(200, "index", nil)
		}

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
