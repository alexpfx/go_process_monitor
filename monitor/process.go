package monitor

import (
	"bufio"
	"github.com/alexpfx/go_process_monitor/util"
	"os/exec"
	"time"
)

type Process struct {
	CmdPath string
	Args    []string

}

func (p Process) Start() chan string {
	wrappedArgs := p.Args
	return start(p.CmdPath, wrappedArgs)
}

func start(cmdPath string, args []string) (ch chan string) {
	ch = make(chan string)

	cmd := exec.Command(cmdPath, args...)

	stdout, err := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	util.Check(err)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			ch <- text
			time.Sleep(5)
		}
	}()

	return
}
