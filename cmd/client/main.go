package main

import (
	"context"
	"fmt"
	"github.com/alexpfx/go_process_monitor/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"io"
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
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() < 2 {
				return fmt.Errorf("argumento faltante. Uso:\n client [OPTIONS...] <pipe_file> <pattern>")
			}

			host := ctx.String("host")
			port := ctx.Int("port")
			url := fmt.Sprintf("%s:%d", host, port)
			allMsgs := ctx.Bool("all_msgs")

			pipeFile := ctx.Args().Get(0)
			pattern := ctx.Args().Get(1)

			conn, err := grpc.Dial(url, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("não pode conectar ao serviço: %s %s", url, err.Error())
			}

			defer conn.Close()

			client := pb.NewProcessMonitorClient(conn)

			stream, err := client.RunOnce(context.Background(), &pb.ConfigOnceRequest{
				PipeName: pipeFile,
				Config: &pb.ConfigChange{
					Pattern:       pattern,
					ReceiveOutput: allMsgs,
				},
			})

			if err != nil {
				return err
			}
			for {
				event, err := stream.Recv()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					return err
				}
				fmt.Println("received event: ")
				fmt.Printf("time %v, text %s, match %v", event.Time, event.Text, event.Match)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
