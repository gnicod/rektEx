package main

import (
	"fmt"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	r "github.com/dancannon/gorethink"
	"github.com/gnicod/rektEx/api"
	"github.com/go-martini/martini"
)

var session *r.Session

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "logs",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "ok bitches")
	})
	m.Group("/api", func(r martini.Router) {
		r.Post("/log", binding.Bind(api.Log{}), api.NewLog)
		r.Get("/log/:appname", api.GetLogForApp)
	})
	m.Run()
}
