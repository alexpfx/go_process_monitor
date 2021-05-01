package main

import (
	"fmt"
	"github.com/alexpfx/go_process_monitor/internal/client"
	"github.com/uMrfave/cli/v2"
	"log"
	"os"
)

var host string
var port int
var pattern string

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       50051,
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "host",
				Aliases:     []string{"u"},
				Value:       "0.0.0.0",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "pattern",
				Aliases:     []string{"s"},
				Required:    true,
				Destination: &pattern,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("usage:\n ps cmd [args...]")
			}
			cmd := c.Args().First()
			args := c.Args().Tail()

			err := client.StartAndListen(host, port, cmd, args, pattern)
			return err
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
