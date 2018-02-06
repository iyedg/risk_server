package main

import (
	"log"
	"net/http"
)

func main() {

	hub := newHub()
	http.Handle("/ws", hub)
	http.Handle("/", &templateHandler{filename: "client.html"})
	go hub.run()
	log.Printf("Listening now on http://%s:80", localIP())
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
