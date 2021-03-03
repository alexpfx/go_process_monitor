package monitor

import (
	"fmt"
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
	ServerChan  chan Msg
	ProcessChan chan string

	//map[term][]process
	queue map[string][]string
	page  int
	cache []string
}

func (m *Monitor) Start() {
	m.queue = make(map[string][]string)
	m.cache = make([]string, 0)
	for {
		m.monitoring()
	}
}

func (m *Monitor) monitoring() {

	select {
	case ms := <-m.ServerChan:
		if ms.Command == Register {
			m.register(ms.Process, ms.Term)
		} else {
			m.unregister(ms.Process, ms.Term)
		}
	case line := <-m.ProcessChan:
		m.page++
		m.cache = append(m.cache, line)
		if m.page%20 == 0 {
			for _, l := range m.cache {
				fmt.Printf("%s\n", l)
			}
		}

		pss, t := m.search(line)
		spawnProcess(pss, line, t)
	}

}

func spawnProcess(pss []string, line string, t string) {
	for _, s := range pss {
		fmt.Printf("Ouvindo termo: %s ps: %s", fmt.Sprintf(InfoColor, t), fmt.Sprintf(NoticeColor, s))
	}
}

func (m Monitor) search(line string) (pss []string, t string) {
	if len(m.queue) == 0 {
		return
	}
	for term := range m.queue {
		if strings.Contains(line, term) {
			pss = m.queue[term]
			t = term
			return
		}
	}
	return
}

//registra um processo que observa um termo
func (m *Monitor) register(ps string, term string) {
	pss := m.queue[term]

	for _, p := range pss {
		if p == ps {
			return
		}
	}
	pss = append(pss, ps)
	m.queue[term] = pss
}

func (m *Monitor) unregister(ps, term string) {
	if _, ok := m.queue[term]; ok {
		delete(m.queue, ps)
	}
}
