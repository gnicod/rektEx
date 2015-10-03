package main

import (
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

func socketHandler(w http.ResponseWriter, r *http.Request) {
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
}
