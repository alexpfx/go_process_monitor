package main

import (
	"fmt"

	"github.com/alexpfx/go_process_monitor/internal/psmon"

	"github.com/urfave/cli/v2"

	"log"

	"os"
	"path/filepath"
)

var host string
var port int
var dir string

func main() {
	homeDir, _ := os.UserHomeDir()
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Aliases:     []string{"u"},
				Value:       "0.0.0.0",
				Destination: &host,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       50051,
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"d"},
				Value:       filepath.Join(homeDir, ".cache/"),
				Destination: &dir,
			},
		},
		Action: func(ctx *cli.Context) error {
			server := psmon.NewServer()

			host := ctx.String("host")
			port := ctx.Int("port")
			err := server.Start(fmt.Sprintf("%s:%d", host, port))

			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
