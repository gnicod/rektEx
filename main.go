package main

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

func main() {
	//go ServeWebsocket()
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
	m.Use(martini.Static("assets"))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "index")
	})
	// api group
	m.Group("/api", func(r martini.Router) {
		r.Post("/log", binding.Bind(Log{}), NewLog)
		r.Get("/log/:appname", GetLogForApp)
	})

	//websocket
	m.Get("/sock", socketHandler)
	m.Run()
}
