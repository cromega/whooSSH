package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"io"

	"bufio"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

var (
	incoming chan []byte
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	incoming = make(chan []byte, 100)

	http.HandleFunc("/whooSSH", WSSHandler)
	handleStaticHTTP()

	stop := make(chan bool, 1)

	go func() {
		<-sigs
		fmt.Println("signal caught, quitting")
		stop <- true
	}()

	startHTTPServer()

	input := make(chan string, 100)

	go func() {
		handleIncomingMessages(input)
	}()

	output := make(chan string, 100)
	startProcess(input, output)
	go func() {
		for line := range output {
			fmt.Print(line)

		}
	}()

	<-stop
	close(incoming)
}

func startProcess(input, output chan string) {
	cmd := exec.Command("bash")
	cmd.Env = os.Environ()

	handle, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	go func() {
		for i := range input {
			fmt.Println("sending data to bash: ", i)

			_, err := io.WriteString(handle, i)
			if err != nil {
				fmt.Println("write failed", err)
			}
		}
	}()

	go func() {
		r := bufio.NewScanner(handle)
		r.Split(bufio.ScanBytes)

		for r.Scan() {
			output <- r.Text()
		}
		fmt.Println("end scan")
	}()
}

func handleStaticHTTP() {
	dir := http.FileServer(http.Dir("./public"))
	http.Handle("/", dir)
}

func startHTTPServer() {
	go func() {
		fmt.Println("Starting Static HTTP server")
		panic(http.ListenAndServe(":8080", nil))
	}()
}

func handleIncomingMessages(input chan string) {
	line := ""
	for message := range incoming {
		if message[0] == 13 {
			fmt.Println("message received: ", line)
			input <- line + "\n"
			line = ""
		} else {
			msg := string(message)
			line += msg
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WSSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

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
