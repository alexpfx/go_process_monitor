package client

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"io"
	"strings"
)

func StartAndListen(host string, port int, cmd string, args []string, pattern string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := pb.NewProcessMonitorClient(conn)

	stream, err := client.StartProcess(context.Background())
	if err != nil {
		return err
	}
	cfgMsg := createPsRequest(cmd, args)
	err = stream.Send(cfgMsg)
	if err != nil {
		return err
	}
	initialFilter := createFilterRequest(pattern, true)
	err = stream.Send(initialFilter)
	if err != nil {
		return err
	}

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil{
			return err
		}
		fmt.Println(rec)
	}
	return nil
}

func createFilterRequest(pattern string, recvAll bool) *pb.ProcessMonitorRequest {
	return &pb.ProcessMonitorRequest{
		Payload: &pb.ProcessMonitorRequest_Filter{
			Filter: &pb.Filter{
				Pattern:       pattern,
				ReceiveOutput: recvAll,
			},
		},
	}
}

func createPsRequest(cmd string, args []string) *pb.ProcessMonitorRequest {
	return &pb.ProcessMonitorRequest{
		Payload: &pb.ProcessMonitorRequest_Process{
			Process: &pb.Process{
				CmdPath: cmd,
				Args:    strings.Join(args, " "),
			},
		},
	}
}
