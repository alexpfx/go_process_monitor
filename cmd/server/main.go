package main

import (
	"github.com/alexpfx/go_process_monitor/monitor"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   50051,
			},
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"u"},
				Value:   "0.0.0.0",
			},
			&cli.StringFlag{Name: "dir", Aliases: []string{"d"},
				Value: filepath.Join(homeDir, ".cache")},
		},
		Action: func(ctx *cli.Context) error {
			port := ctx.Int("port")
			host := ctx.String("host")
			dir := ctx.String("dir")

			runServer(host, port, dir)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(host string, port int, dir string) {
	server := monitor.NewServer(host, port, dir)
	server.Start()
}
