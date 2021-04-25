package monitor

import (
	"fmt"
	"github.com/alexpfx/go_process_monitor/util"
	"os/exec"
	"regexp"
	"strings"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%v\033[0m"
	WarningColor = "\033[1;33m%v\033[0m"
	ErrorColor   = "\033[1;31m%v\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

type Monitor struct {
	queue map[string][]Observer
	page  int
	cache []string
}

func (m *Monitor) Start(srvCh chan Observer, prCh chan string) {
	m.queue = make(map[string][]Observer)
	m.cache = make([]string, 0)
	m.monitoring(srvCh, prCh)
}

func (m *Monitor) monitoring(srvCh chan Observer, psCh chan string) {
	for {
		select {
		case ms := <-srvCh:
			switch ms.Command {
			case Register:
				m.register(ms)
			case Unregister:
				m.unregister(ms)
			}
		case line, ok := <-psCh:
			pss, _ := m.search(line)

			fmt.Printf("%.70s ...", line)
			spawnProcess(pss, line)
			fmt.Println()
			if !ok {
				psCh = nil
			}
		}

		if psCh == nil {
			break
		}

	}
}

func spawnProcess(pss []Observer, line string) {
	if len(pss) == 0 {
		return
	}
	for _, msg := range pss {

		var extractArgs []string
		if msg.Pattern != "" {
			regex := regexp.MustCompile(msg.Pattern)
			extractArgs = regex.FindAllString(line, -1)
		}

		split := strings.Split(msg.Process, " ")
		cmdName := split[0]
		cmdArgs := split[1:]

		cmd := exec.Command(cmdName, append(cmdArgs, extractArgs...)...)

		output, err := cmd.CombinedOutput()
		util.Check(err)

		fmt.Printf("%70s", "-> "+string(output))
	}

}

func (m Monitor) search(line string) (pss []Observer, found string) {
	if len(m.queue) == 0 {
		return
	}

	for term := range m.queue {
		if strings.Contains(line, term) {
			pss = m.queue[term]
			found = term

			break
		}
	}
	return
}

//registra um processo que observa um termo
func (m *Monitor) register(msg Observer) {
	msgs := m.queue[msg.Term]

	for _, m := range msgs {
		if m.Process == msg.Process {
			return
		}
	}
	msgs = append(msgs, msg)
	m.queue[msg.Term] = msgs
}

func (m *Monitor) unregister(msg Observer) {
	if _, ok := m.queue[msg.Term]; ok {
		delete(m.queue, msg.Term)
	}
}
