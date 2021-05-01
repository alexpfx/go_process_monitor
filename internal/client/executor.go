package client

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"strings"
)

func StartAndListen(host string, port int, cmd string, args []string, pattern string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewProcessMonitorClient(conn)

	stream, err := client.StartProcess(context.Background())
	if err != nil {
		return err
	}

	waitChan := make(chan struct{})
	go func() {
		for {
			log.Println("iniciando escuta no stream")
			_, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF: recebimento encerrado")
				close(waitChan)
				break
			}
			if err != nil {
				log.Println("erro no Recv ", err.Error())
				break
			}
		}
	}()


	cfgMsg := createPsRequest(cmd, args)
	_ = stream.Send(cfgMsg)
	initialFilter := createFilterRequest(pattern, true)
	_ = stream.Send(initialFilter)
	_ = stream.CloseSend()

	<- waitChan
	return nil
}

func createFilterRequest(pattern string, recvAll bool) *pb.ProcessMonitorRequest {
	return &pb.ProcessMonitorRequest{
		Msg: &pb.ProcessMonitorRequest_Filter{
			Filter: &pb.Filter{
				Pattern:       pattern,
				ReceiveOutput: recvAll,
			},
		},
	}
}

func createPsRequest(cmd string, args []string) *pb.ProcessMonitorRequest {
	return &pb.ProcessMonitorRequest{
		Msg: &pb.ProcessMonitorRequest_Process{
			Process: &pb.Process{
				CmdPath: cmd,
				Args:    strings.Join(args, " "),
			},
		},
	}
}
