package psmon

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
)

type server struct {
	mutex   sync.RWMutex
	psStore ProcessStore
}

func (s *server) Start(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProcessMonitorServer(grpcServer, &server{
		mutex:   sync.RWMutex{},
		psStore: NewProcessStore(),
	})

	log.Printf("servidor ouvindo em %s", address)
	err = grpcServer.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
func (s *server) RunProcess(ctx context.Context, rq *pb.RunRequest) (*pb.RunResponse, error) {
	uid := rq.GetPsUid()
	ps := s.psStore.Get(uid)
	if ps == nil {
		return nil, fmt.Errorf("Process with uid %s not found", uid)
	}

	rp := &pb.RunResponse{
		Process: &pb.Process{
			Name: ps.Path,
			Args: strings.Join(ps.Args, " "),
		},
	}
	//exec
	log.Println("running...", rp)

	return rp, nil
}
func (s *server) RegisterProcess(ctx context.Context, rq *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	process := rq.Process.GetName()
	args := rq.Process.GetArgs()

	spArgs := strings.Split(args, " ")
	cmd := exec.Command(process, spArgs...)

	uid := s.psStore.Save(cmd)
	log.Printf("processo registrado (%d). uid: %s\n", s.psStore.Count(), uid)
	return &pb.RegisterResponse{
		PsUid: uid,
	}, nil
}
