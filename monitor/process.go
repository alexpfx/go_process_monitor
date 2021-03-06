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

func (p Process) Start(ch chan string)  {

	cmd := exec.Command(p.CmdPath, p.Args...)

	stdout, err := cmd.StdoutPipe()
	util.Check(err)
	scanner := bufio.NewScanner(stdout)

	util.Check(err)

	err = cmd.Start()
	util.Check(err)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			ch <- text
		}
		time.Sleep(time.Millisecond * 300)
		close(ch)
	}()

	err = cmd.Wait()
	util.Check(err)
}

