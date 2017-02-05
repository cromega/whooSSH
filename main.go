package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	incoming chan []byte
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	incoming = make(chan []byte, 100)

	stop := make(chan bool, 1)

	go func() {
		<-sigs
		fmt.Println("signal caught, quitting")
		stop <- true
	}()

	NewServer()

	sp, err := NewSubProcess("bash")
	if err != nil {
		panic(err)
	}
	defer sp.kill()

	go func() {
		handleIncomingMessages(sp.input)
	}()

	go func() {
		for line := range sp.output {
			fmt.Print(line)
		}
	}()

	<-stop
	close(incoming)
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
