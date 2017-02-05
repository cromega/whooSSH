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

	<-stop
	close(incoming)
}
