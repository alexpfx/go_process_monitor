package monitor

import (
	"bufio"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"
)

type Server interface {
	Start()
}

type server struct {
	host string
	port int
}

func (s *server) StartProcess(req *pb.ProcessMonitorRequest, stream pb.ProcessMonitor_StartProcessServer) error {
	pattern := req.GetPattern()
	_ = req.GetReceiveOutput()

	cmd := exec.Command(req.GetCmdPath(), req.GetArgs())
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	scanner := bufio.NewScanner(pr)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()

			if strings.Contains(text, pattern) {
				err := stream.Send(&pb.Event{
					Time:  time.Now().Unix(),
					Text:  text,
					Match: true,
				})
				if err != nil {
					return
				}
			}
		}
	}()

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func RunServer(host string, port int) error {
	address := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("não pode iniciar o servidor [%s]: %v",
			address, err)

	}
	grpcServer := grpc.NewServer()
	pb.RegisterProcessMonitorServer(grpcServer, &server{})
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("não pode iniciar o servidor [%s]: %v",
			address, err)
	}
	return nil
}
