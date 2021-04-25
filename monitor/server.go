package monitor

import (
	"bufio"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type server struct {
	address string
	port    int
	dir     string
}

func (s *server) RunOnce(request *pb.ConfigOnceRequest, stream pb.ProcessMonitor_RunOnceServer) error {
	pipeName := request.PipeName
	cfg := request.Config

	if p, ok := peer.FromContext(stream.Context()); ok {
		log.Printf("cliente conectado: %s", p.Addr)
	}

	ch, err := openPipe(s.dir, pipeName)
	if err != nil{
		return err
	}

	for msg := range ch{
		if ev, shouldSend := buildEvent(cfg, msg); shouldSend{
			_ = stream.Send(ev)
		}
	}
	return nil
}

func (s *server) Run(stream pb.ProcessMonitor_RunServer) error {
	if p, ok := peer.FromContext(stream.Context()); ok {
		log.Printf("cliente conectado: %s", p.Addr)
	}
	recv, err := stream.Recv()
	if err != nil {
		return err
	}
	pipeName := recv.GetPipeName()
	ch, err := openPipe(s.dir, pipeName)
	if err != nil {
		return err
	}
	changes := configChanges(stream)

	cfg := pb.ConfigChange{}
	select {
	case msg := <-ch:
		if ev, shouldSend := buildEvent(&cfg, msg); shouldSend {
			_ = stream.Send(ev)
		}
	case ncf := <-changes:
		cfg = pb.ConfigChange{
			Pattern:       ncf.Pattern,
			ReceiveOutput: ncf.ReceiveOutput,
		}
	}
	return nil
}

func buildEvent(cfg *pb.ConfigChange, msg string) (*pb.Event, bool) {
	pattern := cfg.Pattern
	rcvOut := cfg.ReceiveOutput
	if pattern == "" {
		return nil, false
	}
	ev := &pb.Event{
		Time:  time.Now().Unix(),
		Text:  msg,
		Match: strings.Contains(msg, pattern),
	}

	if rcvOut {
		return ev, true
	}
	return ev, ev.Match
}

func configChanges(stream pb.ProcessMonitor_RunServer) chan *pb.ConfigChange {
	ch := make(chan *pb.ConfigChange)
	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				close(ch)
				break
			}
			if err != nil {
				close(ch)
				log.Fatal(err)
			}
			newCfg := recv.GetConfig()
			ch <- &pb.ConfigChange{
				Pattern:       newCfg.Pattern,
				ReceiveOutput: newCfg.ReceiveOutput,
			}
		}
	}()
	return ch
}

func openPipe(dir string, name string) (chan string, error) {
	ch := make(chan string)
	filePath := filepath.Join(dir, name)
	pipeFile, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(pipeFile)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			ch <- text
		}
	}()

	return ch, nil
}

func NewServer(address string, port int, dir string) Server {
	return &server{
		address: address,
		port:    port,
		dir:     dir,
	}
}

func (s *server) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		log.Fatalf("Não pode iniciar o servidor [%s:%d]: %v",
			s.address, s.port, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProcessMonitorServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Não pode iniciar o servidor [%s:%d]: %v",
			s.address, s.port, err)
	}

}

type Server interface {
	Start()
}
