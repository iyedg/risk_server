package main

import "github.com/gorilla/websocket"

type client struct {
	id     string
	socket *websocket.Conn
	hub    *hub
	send   chan []byte
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.hub.broadcast <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func newClient(socket *websocket.Conn, hub *hub) *client {
	return &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		hub:    hub,
	}
}
