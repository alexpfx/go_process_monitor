package main

import (
	"fmt"
	"github.com/alexpfx/go_process_monitor/monitor"
	"github.com/alexpfx/go_process_monitor/util"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%v\033[0m"
	WarningColor = "\033[1;33m%v\033[0m"
	ErrorColor   = "\033[1;31m%v\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var eventCmd = make([]string, 0)
var template = "%s"
var pattern *regexp.Regexp
var wg sync.WaitGroup

func main() {
	app := cli.App{

		Commands: []*cli.Command{
			{
				Name: "monitor",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "porta onde deve esperar conexões http",
						Value:   7003,
					},
					&cli.StringFlag{
						Name:     "process",
						Aliases:  []string{"s"},
						Usage:    "processo que terá sua saída monitorada",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "auth",
						Aliases: []string{"a"},
						Usage:   "user:pass",

						Required: true,
					},
				},

				Action: func(c *cli.Context) error {

					auth := strings.Split(c.String("auth"), ":")
					if len(auth) != 2 {
						return fmt.Errorf("user e pass são obrigatórios")
					}

					psSplitted := strings.Split(c.String("s"), " ")
					cmdPath := psSplitted[0]
					cmdArgs := psSplitted[1:]

					ps := monitor.Process{
						CmdPath: cmdPath,
						Args:    cmdArgs,
					}

					psCh := make(chan string)
					go func() {
						ps.Start(psCh)
					}()

					srv := monitor.Server{
						Port: c.Int("port"),
						User: auth[0],
						Pass: auth[1],
					}

					srvCh := make(chan monitor.Msg)
					go func() {
						srv.Start(srvCh)
					}()

					mon := monitor.Monitor{}

					mon.Start(srvCh, psCh)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	util.Check(err)

}

func search(queue map[string]bool, text string) string {
	for k, v := range queue {
		if !v {
			continue
		}

		sprintf := fmt.Sprintf(template, k)
		if strings.Contains(text, sprintf) {
			return k
		}
	}
	return ""
}

/*
func countQueue() int {
	i := 0
	for _, v := range queue {
		if v {
			i++
		}
	}
	return i
}
*/
//func printQueue(writer http.ResponseWriter) {
//	writer.WriteHeader(http.StatusOK)
//
//	fmt.Fprintf(writer, "fila: \n")
//	for k, v := range queue {
//		if v {
//			fmt.Fprintln(writer, k)
//		}
//	}
//}

func notify(text string, str string) {
	fmt.Printf(InfoColor, str+": ")
	fmt.Printf(DebugColor, text+"\n")

	/*if eventCmd != nil {
		cmd := exec.Command(eventCmd[0], append(eventCmd, str...)...)

		out, err := cmd.CombinedOutput()
		Check(err)
		fmt.Printf(DebugColor, string(out))
	}*/

}
