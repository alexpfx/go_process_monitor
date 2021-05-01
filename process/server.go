package process

import (
	"fmt"
	"log"
	"net"
	"path"
)

type Monitor struct {
	Dir      string
	Filename string
}

func (m Monitor) RunAndServ(){
	tempFile := path.Join(m.Dir, m.Filename)
	listen, err := net.Listen("unix", tempFile)
	log.Println(tempFile)
	defer listen.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil{
			log.Fatal(err)
		}
		fmt.Println(conn)
		fmt.Println("accepted")

	}

}


