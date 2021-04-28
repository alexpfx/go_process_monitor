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

func (s *server) StartProcess(stream pb.ProcessMonitor_StartProcessServer) error {
	psChan, filterChan, err := splitRecv(stream)
	if err != nil {
		return err
	}
	ps := <-psChan

	cmd := exec.Command(ps.CmdPath, ps.Args)
	ch, qch, err := handleCmd(cmd, filterChan)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-ch:
			fmt.Println(msg)
		case q := <-qch:
			fmt.Println(q)
		default:
		}
	}
	return nil
}

func handleCmd(cmd *exec.Cmd, filter chan *pb.Filter) (chan *pb.Event, chan bool, error) {
	ch := make(chan *pb.Event)
	quitCh := make(chan bool)

	r, w := io.Pipe()
	cmd.Stdout = w
	scanner := bufio.NewScanner(r)
	var f *pb.Filter

	go func() {
		for {
			ff := <-filter
			f = ff
			fmt.Println(f)
		}
	}()
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if f == nil {
				continue
			}
			fmt.Println(text)
			if strings.Contains(text, f.Pattern) {
				ch <- &pb.Event{
					Time:  time.Now().Unix(),
					Text:  text,
					Match: true,
				}
				return
			}
			if !f.ReceiveOutput {
				continue
			}
			ch <- &pb.Event{
				Time:  time.Now().Unix(),
				Text:  text,
				Match: true,
			}
			return
		}
	}()
	return ch, quitCh, nil
}

func splitRecv(stream pb.ProcessMonitor_StartProcessServer) (chan *pb.Process, chan *pb.Filter, error) {
	chp := make(chan *pb.Process)
	chf := make(chan *pb.Filter)

	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				close(chp)
				close(chf)
			}
			if err != nil {
				return
			}

			switch recv.Payload.(type) {
			case *pb.ProcessMonitorRequest_Process:
				chp <- recv.GetProcess()
			case *pb.ProcessMonitorRequest_Filter:
				chf <- recv.GetFilter()
			}
		}
	}()
	return chp, chf, nil
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
