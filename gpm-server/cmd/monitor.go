package main

import (
	"github.com/alexpfx/go_process_monitor/process"
	"io/ioutil"
)

func main() {
	dir, _ := ioutil.TempDir("/tmp", "pm_test*")
	m := process.Monitor{
		Dir:      dir,
		Filename: "test.socket",
	}
	m.RunAndServ()
}