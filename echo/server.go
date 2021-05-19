package echo

import (
	"log"
	"net"
)

type Server interface {
	Run() error
}

func NewServer(socket string) Server {
	return server{socket: socket}
}

type server struct {
	socket string
}

func (s server) Run() error {
	listen, err := net.Listen("unix", s.socket)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			err = echoServer(conn)
			if err != nil {
				log.Fatal(err)
			}

		}
	}()

	return nil
}

func echoServer(conn net.Conn) error {
	for {
		buf := make([]byte, 512)
		count, err := conn.Read(buf)
		if err != nil {
			return err
		}
		data := buf[0:count]
		count, err = conn.Write(data)
		if err != nil {
			return err
		}
	}

}
