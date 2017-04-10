package main

import (
	"fmt"
)

type session struct {
	sp       *subProcess
	conn     *connection
	messages chan string
}

func NewSession(conn *connection) (s *session) {
	s = &session{
		conn: conn,
	}

	sp := NewSubProcess("bash")
	s.sp = sp

	return
}

func (s *session) start() (err error) {
	s.sendMessage("hilo")

	err = s.sp.start()
	if err != nil {
		return
	}

	go func() {
		line := ""
		for {
			message := s.conn.Read()
			if len(message) == 0 {
				fmt.Println("connection closed, terminating subprocess")
				s.sp.kill()
				break
			}

			if message[0] == 13 {
				fmt.Println("message received: ", line)
				s.sp.input <- line + "\n"
				line = ""
			} else {
				msg := string(message)
				line += msg
			}
		}
	}()

	go func() {
		for line := range s.sp.output {
			fmt.Print(line)
			s.sendMessage(line)
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
