package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type server struct {
}

type connection struct {
	conn     *websocket.Conn
	messages chan string
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

	session := NewSession(&connection{conn: conn})
	err = session.start()
	if err != nil {
		fmt.Println("session wtfed out: %v", err)
		return
	}

}

func (c *connection) Write(message string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *connection) Read() string {
	messageType, data, err := c.conn.ReadMessage()

	if err != nil {
		fmt.Println(err)
	}

	if messageType == websocket.TextMessage {
		return string(data)
	}

	return ""
}

func handleStaticHTTP() {
	dir := http.FileServer(http.Dir("./public"))
	http.Handle("/", dir)
}
