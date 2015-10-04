package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func socketHandler(args martini.Params, w http.ResponseWriter, r *http.Request) {
	//When a socket client register to an app, listen to this app and automatically send new update to this client
	onChange()
	appname := args["appname"]
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

	logs, _ := GetLogForApp(appname)
	log.Println(logs)
	jsonlogs, err := json.Marshal(logs)
	if err != nil {
		log.Println(err)
	}

	if err := sockCli.websocket.WriteMessage(1, jsonlogs); err != nil {
		return
	}

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
}
