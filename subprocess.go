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

func NewSubProcess(command string) (*subProcess, error) {
	sp := &subProcess{
		input:  make(chan string, 100),
		output: make(chan string, 100),
	}

	cmd := exec.Command(command)
	cmd.Env = os.Environ()

	handle, err := pty.Start(cmd)
	if err != nil {
		return nil, err
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
		r.Split(bufio.ScanBytes)

		for r.Scan() {
			sp.output <- r.Text()
		}
		fmt.Println("end scan")
	}()

	sp.cmd = cmd
	return sp, nil
}

func (sp *subProcess) kill() {
	sp.cmd.Process.Kill()
}
