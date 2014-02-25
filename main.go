package main

import (
	"crypto/sha1"
	"fmt"

	"github.com/codegangsta/martini"
)

func main() {
	m:= martini.Classic()
	// add an auth step here to redirect to /login if not authed
	m.Get("/", func() string {
		return "This is the main page"
	})

	m.Get("/poll/:id", func(params martini.Params) string {
		return "This is the page for poll #" + params["id"]
	})

	m.Get("/hash", func() string {
		h := sha1.New()
		h.Write([]byte("this is a hash"))
		hs := h.Sum(nil)
		return fmt.Sprintf("%x",hs)
	})

	m.Get("/login", func() string {
		return "This is the login page."
	})
	m.Run()
}
