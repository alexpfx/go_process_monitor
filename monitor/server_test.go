package monitor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path"
	"testing"
	"time"
)

func Test_server_Subscribe(t *testing.T) {
	tempDir, err := ioutil.TempDir("/tmp/", "pm")
	if err != nil {
		log.Fatal(err)
	}

	tempFile := path.Join(tempDir, "sock.file")
	listen, err := net.Listen("unix", tempFile)
	defer listen.Close()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			all, err := ioutil.ReadAll(conn)
			if err != nil {
				return
			}
			fmt.Println(all)
		}
	}()

	time.Sleep(time.Second)
	dial, err := net.Dial("unix", tempFile)

	writer := bufio.NewWriter(dial)
	writer.WriteString("Hello world")
	if err != nil {
		log.Fatal(err)
	}

	select {

	}
}
