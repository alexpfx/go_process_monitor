package process

import (
	"context"
	"github.com/alexpfx/go_process_monitor/pb"
)

type server struct {

}



func (s server) RegisterProcess(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	panic("implement me")
}

func (s server) Subscribe(request *pb.SubscribeRequest, subscribeServer pb.ProcessMonitor_SubscribeServer) error {
	panic("implement me")
}

