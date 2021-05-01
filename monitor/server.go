package monitor

import (
	"bufio"
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type server struct {
	process map[uint64]*exec.Cmd
}

func (s *server) Subscribe(req *pb.SubscribeRequest, stream pb.ProcessMonitor_SubscribeServer) error {
	cmd := s.process[req.KeyId]
	if cmd == nil || cmd.ProcessState == nil {
		return fmt.Errorf("processo nao encontrado")
	}
	return nil
}

func (s *server) PrepareProcess(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	if len(s.process) == 0 {
		s.process = make(map[uint64]*exec.Cmd)
	}

	ps := req.GetProcess()
	if ps == nil {
		return nil, fmt.Errorf("deveria passar o Processo")
	}

	rand.Seed(time.Now().Unix())
	key := rand.Uint64()

	s.process[key] = exec.Command(ps.GetCmdPath(), ps.GetArgs())
	return &pb.ProcessResponse{
		KeyId: key,
	}, nil
}

func (s *server) StartProcess(stream pb.ProcessMonitor_StartProcessServer) error {

	psChan, filterChan := splitRecv(stream)
	ps, filter := waitCfg(psChan, filterChan)
	cmd := exec.Command(ps.CmdPath, ps.Args)
	outputChan, err := handleCmd(cmd, filter)

	if err != nil {
		return err
	}
	log.Println("iniciando loop de escuta do processo")

	for {
		msg := <-outputChan
		err = stream.Send(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func waitCfg(psChan chan *pb.Process, filterChan chan *pb.Filter) (*pb.Process, *pb.Filter) {
	return <-psChan, <-filterChan
}

func handleCmd(cmd *exec.Cmd, filter *pb.Filter) (chan *pb.Event, error) {
	eventChan := make(chan *pb.Event)
	r, w := io.Pipe()

	cmd.Stdout = w
	scanner := bufio.NewScanner(r)

	fmt.Println("gofunc 1")
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			e := pb.Event{
				Time: time.Now().Unix(),
				Text: text,
			}

			if strings.Contains(text, filter.Pattern) {
				e.Match = true
				eventChan <- &e
			} else if filter.ReceiveOutput {
				eventChan <- &e
			}
		}
		fmt.Println("nada mais a ler")
		close(eventChan)
		eventChan = nil
	}()
	fmt.Println("gofunc 2")
	go func() {
		log.Println("iniciando commando")
		err := cmd.Run()
		if err != nil {
			log.Println("erro ao iniciar o comando")
			return
		}
		log.Println("comando parou")
		err = r.Close()
		if err != nil {
			fmt.Println("erro ao fechar reader")

		}
		err = w.Close()
		if err != nil {
			fmt.Println("erro ao fechar writer")
		}
	}()
	return eventChan, nil
}

func splitRecv(stream pb.ProcessMonitor_StartProcessServer) (chan *pb.Process, chan *pb.Filter) {
	psChan := make(chan *pb.Process)
	filterChan := make(chan *pb.Filter)

	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				close(psChan)
				close(filterChan)
				break
			}
			if err != nil {
				log.Fatal("error no recv ", err)
				return
			}

			switch recv.Msg.(type) {
			case *pb.ProcessMonitorRequest_Process:
				psChan <- recv.GetProcess()
			case *pb.ProcessMonitorRequest_Filter:
				filterChan <- recv.GetFilter()
			}
		}
	}()
	return psChan, filterChan
}
