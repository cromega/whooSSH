package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type server struct {
	conn *websocket.Conn
}

func NewServer() *server {
	server := &server{}
	http.HandleFunc("/whooSSH", server.WSSHandler)
	handleStaticHTTP()

	go func() {
		fmt.Println("Starting HTTP server")
		panic(http.ListenAndServe(":8080", nil))
	}()

	return server
}

func (s *server) WSSHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.conn = conn

	conn.WriteMessage(websocket.TextMessage, []byte("lo!"))

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}

		if messageType == websocket.TextMessage {
			incoming <- data
		}
	}
}

func handleStaticHTTP() {
	dir := http.FileServer(http.Dir("./public"))
	http.Handle("/", dir)
}
