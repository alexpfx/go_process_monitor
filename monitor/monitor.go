package monitor

import (
	"fmt"
	"strings"
	"time"
)

type Monitor struct {
	ServerChan  chan Msg
	ProcessChan chan string

	//map[term][]process
	queue map[string][]string
}

func (m *Monitor) Start() {
	m.queue = make(map[string][]string)
}

func (m Monitor) monitoring() {
	for {
		if len(m.queue) == 0 {
			time.Sleep(10)
			continue
		}
		select {
		case ms := <-m.ServerChan:
			if ms.Command == Register {
				m.register(ms.Process, ms.Term)
			} else {
				m.unregister(ms.Process, ms.Term)
			}
		case line := <-m.ProcessChan:
			pss := m.search(line)
			spawnProcess(pss, line)
		}
	}
}

func spawnProcess(pss []string, line string) {
	for i, s := range pss {
		fmt.Printf("Ouvindo termo: %d %s", i, s)
	}
}

func (m Monitor) search(line string) (pss []string) {
	if len(m.queue) == 0 {
		return
	}
	for term := range m.queue {
		if strings.Contains(line, term) {
			pss = m.queue[term]
			return
		}
	}
	return
}

//registra um processo que observa um termo
func (m *Monitor) register(ps string, term string) {
	pss, exist := m.queue[term]

	if !exist {
		pss = append(pss, ps)
		return
	}

	for _, p := range pss {
		if p == ps {
			return
		}
	}
	pss = append(pss, ps)

}

func (m *Monitor) unregister(ps, term string) {
	if _, ok := m.queue[term]; ok {
		delete(m.queue, ps)
	}
}
