package process

import (
	"log"
	"net"
	"testing"
)

func TestMonitor_RunAndServ_connection(t *testing.T) {

	dial, _ := net.Dial("unix", "/tmp/pm_test846665418/test.socket")
	count, err := dial.Write([]byte("Hello world"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("writed ", count)

}
