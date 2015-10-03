package main

import (
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	m.Get("/sock", func(w http.ResponseWriter, r *http.Request) {
		log.Println(ActiveClients)
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		client := ws.RemoteAddr()
		sockCli := ClientConn{ws, client}
		addClient(sockCli)

		for {
			log.Println(len(ActiveClients), ActiveClients)
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				deleteClient(sockCli)
				log.Println("bye")
				log.Println(err)
				return
			}
			broadcastMessage(messageType, p)
		}
	})
	m.Run()
}
