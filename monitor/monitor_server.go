package monitor

import (
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	Host string
	Port int
}

func (s *Server) Start() error{
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("não pode iniciar o servidor [%s]: %v",
			address, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterProcessMonitorServer(grpcServer, &server{})
	log.Println("server running in ", address)
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("não pode iniciar o servidor [%s]: %v",
			address, err)
	}
	return nil
}
