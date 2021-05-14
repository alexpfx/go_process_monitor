package psmon

import "os/exec"

//import "math/rand"
//import "time"
import "github.com/google/uuid"
import "fmt"

func NewProcessStore() ProcessStore {
	p := new(psStore)
	p.store = make(map[string]*exec.Cmd)
	return p

}

type ObserverStore interface {
}

type ProcessStore interface {
	Save(cmd *exec.Cmd) string
	Get(uid string) *exec.Cmd
	Count() int
	Remove(uid string)
}

type psStore struct {
	store map[string]*exec.Cmd
}

func (p psStore) Get(uid string) *exec.Cmd {
	return p.store[uid]
}

func (p *psStore) Save(cmd *exec.Cmd) string {
	uid := fmt.Sprintf("%s", uuid.New())
	p.store[uid] = cmd
	return uid
}
func (p psStore) Count() int {
	return len(p.store)
}
func (p *psStore) Remove(uid string) {
	delete(p.store, uid)
}
