package process

import (
	"bufio"
	"github.com/alexpfx/go_process_monitor/internal/psmon"
	"log"
	"os/exec"
)

type monitor struct {
	cmd *exec.Cmd
}

func NewMonitor(cmd *exec.Cmd) psmon.Source {
	return monitor{
		cmd: cmd,
	}
}

func (m monitor) Generate() (chan string, error) {
	return startCmd(m.cmd)
}

func startCmd(cmd *exec.Cmd) (chan string, error) {
	ch := make(chan string)

	outPipe, _ := cmd.StdoutPipe()

	scanner := bufio.NewScanner(outPipe)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			ch <- text
		}
	}()
	go func() {
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		close(ch)
	}()
	return ch, nil
}
