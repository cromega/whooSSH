package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

type subProcess struct {
	input  chan string
	output chan string
	cmd    *exec.Cmd
}

func NewSubProcess(command string) *subProcess {
	sp := &subProcess{
		input:  make(chan string, 100),
		output: make(chan string, 100),
	}

	cmd := exec.Command(command)
	cmd.Env = os.Environ()

	sp.cmd = cmd
	return sp
}

func (sp *subProcess) start() (err error) {
	handle, err := pty.Start(sp.cmd)
	if err != nil {
		return
	}

	go func() {
		for i := range sp.input {
			fmt.Println("sending data to bash: ", i)

			_, err := io.WriteString(handle, i)
			if err != nil {
				fmt.Println("write failed", err)
			}
		}
	}()

	go func() {
		r := bufio.NewScanner(handle)
		r.Split(bufio.ScanRunes)

		for r.Scan() {
			sp.output <- r.Text()
		}
		fmt.Println("end scan")
	}()

	return
}

func (sp *subProcess) kill() {
	sp.cmd.Process.Kill()
}
