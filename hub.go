package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type hub struct {
	join      chan *client
	leave     chan *client
	clients   map[*client]bool
	broadcast chan []byte
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
			log.Printf("%v has joined", client.id)
		case client := <-h.leave:
			defer close(client.send)
			delete(h.clients, client)
			log.Printf("%v has left", client.id)
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func (h *hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := newClient(socket, h)
	name, ok := r.URL.Query()["name"]
	if ok {
		client.id = name[0]
	}
	h.join <- client
	defer func() { h.leave <- client }()
	go client.write()
	client.read()
}

func newHub() *hub {
	return &hub{
		broadcast: make(chan []byte),
		join:      make(chan *client),
		leave:     make(chan *client),
		clients:   make(map[*client]bool),
	}
}
