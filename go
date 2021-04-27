package main

import (
	"fmt"
	"github.com/alexpfx/go_process_monitor/internal/client"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
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
			&cli.BoolFlag{
				Name:    "all_msgs",
				Aliases: []string{"a"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "ps",
				Action: func(c *cli.Context) error {
					host := c.String("host")
					port := c.Int("port")

					if c.NArg() < 1 {
						return fmt.Errorf("usage:\n ps cmd [args...]")
					}
					cmd := c.Args().First()
					args := c.Args().Tail()
					psCli := client.NewExecPs(host, port, cmd, args)

					_, err := psCli.Run()

					return err
				},
			},
			{
				Name: "subscribe",
				Action: func(c *cli.Context) error {
					fmt.Println("subscribe")
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() < 2 {
				return cli.ShowAppHelp(ctx)
				//return fmt.Errorf("argumento faltante. Uso:\n client [OPTIONS...] <pipe_file> <pattern>")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
