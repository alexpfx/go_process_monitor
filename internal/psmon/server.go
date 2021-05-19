package psmon

import (
	"sync"
)

type Server interface {
	Start(address string) error
}

func NewServer() Server {
	return &server{
		mutex:   sync.RWMutex{},
		psStore: NewProcessStore(),
	}

}
