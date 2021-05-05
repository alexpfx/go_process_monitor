package psmon

import "os/exec"
import "math/rand"
import "time"

import "fmt"

func NewProcessStore() ProcessStore {
	p := new(psStore)
	p.store = make(map[string]*exec.Cmd)
	return p
}

var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type ObserverStore interface {
}

type ProcessStore interface {
	Save(cmd *exec.Cmd) string
}

type psStore struct {
	store map[string]*exec.Cmd
}

func (p *psStore) Save(cmd *exec.Cmd) string {
	uid := fmt.Sprintf("%v", seedRand)
	p.store[uid] = cmd
	return uid
}
