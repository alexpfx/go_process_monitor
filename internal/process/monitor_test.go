package process

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestMonitor_Run(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ggSearch := exec.CommandContext(ctx, "go_search", "-l", "60", "-d", "/media/nv4/CompilationOfManyBreaches", "-p", "castanho")
	source := NewMonitor(ggSearch)

	generate, _ := source.Generate()
	go func() {

		for {
			select {
			case msg := <-generate:
				fmt.Println(msg)
			}

		}
	}()
	time.Sleep(5 * time.Second)

	cancelFunc()
}
