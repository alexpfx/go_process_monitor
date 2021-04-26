package client

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/monitor"
	"github.com/alexpfx/go_process_monitor/pb"
	"google.golang.org/grpc"
	"strings"
)

type execProcess struct {
	host string
	port int
	cmd  string
	args []string
}

func (c execProcess) Run() error {
	host := c.host
	port := c.port

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewProcessMonitorClient(conn)

	res, err := client.ExecProcess(context.Background(), &pb.ExecProcessRequest{
		CmdPath: c.cmd,
		Args:    strings.Join(c.args, " "),
	})
	if err != nil {
		return err
	}

	fmt.Println(res.GetFilePath())
	return nil
}

func NewExecPs(host string, port int, cmd string, args []string) monitor.Client {
	return execProcess{
		host: host,
		port: port,
		cmd: cmd,
		args: args,
	}
}
