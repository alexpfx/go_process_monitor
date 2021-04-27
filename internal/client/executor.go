package client

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"io"
	"strings"
	"time"
)

func StartAndListen(host string, port int, cmd string, args []string, pattern string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewProcessMonitorClient(conn)
	stream, err := client.StartProcess(context.Background(), &pb.ProcessMonitorRequest{
		CmdPath:       cmd,
		Args:          strings.Join(args, " "),
		Pattern:       pattern,
		ReceiveOutput: false,
	})
	if err != nil {
		return err
	}

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("match: [%v %s]\n",
			time.Unix(recv.GetTime(), 0), recv.GetText())

	}
	return nil
}
