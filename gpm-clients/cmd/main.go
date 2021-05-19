package main

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
)

var host string
var port int
var pattern string

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "register",
				Usage: "registra um comando para posterior execução e monitoramento",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "exec",
						Aliases: []string{"-x"},
						Usage:   "executa o comando imediatamente",
						Value:   true,
					},
				},

				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("usage:\n ps cmd [args...]")
					}
					cmd := c.Args().First()
					args := c.Args().Tail()
					fmt.Println(cmd)
					fmt.Println(args)
					host := c.String("host")
					port := c.Int("port")
					conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
					if err != nil {
						return err
					}
					defer conn.Close()
					client := pb.NewProcessMonitorClient(conn)
					req := &pb.RegisterRequest{
						Process: &pb.Process{
							Name: cmd,
							Args: strings.Join(args, " "),
						},
					}
					res, err := client.RegisterProcess(context.Background(), req)
					if err != nil {
						return err
					}
					shouldExec := c.Bool("exec")
					if !shouldExec {
						return nil
					}

					rReq := &pb.RunRequest{
						PsUid: res.PsUid,
					}

					rRes, err := client.RunProcess(context.Background(), rReq)
					log.Println(rRes)

					return nil
				},
			},
		},
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
				Destination: &pattern,
			},
		},
		Action: func(c *cli.Context) error {

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
