package main

import (
	"fmt"
)

type session struct {
	sp       *subProcess
	conn     *connection
	messages chan string
}

func NewSession(conn *connection) (s *session, err error) {
	s = &session{
		conn: conn,
	}

	s.sendMessage("hilo")

	sp := NewSubProcess("bash")
	err = sp.start()
	if err != nil {
		return
	}
	s.sp = sp

	go func() {
		line := ""
		for {
			message := conn.Read()
			if len(message) == 0 {
				fmt.Println("connection closed, terminating subprocess")
				sp.kill()
				break
			}

			if message[0] == 13 {
				fmt.Println("message received: ", line)
				sp.input <- line + "\n"
				line = ""
			} else {
				msg := string(message)
				line += msg
			}
		}
	}()

	go func() {
		for line := range sp.output {
			fmt.Print(line)
			conn.Write(line)
		}
	}()

	return
}

func (s *session) sendMessage(message string) {
	s.conn.Write(message)
}

func (s *session) End() {
	s.sp.kill()
}
