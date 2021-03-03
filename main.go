package main

import (
	"bufio"
	"fmt"
	"github.com/alexpfx/go_process_monitor/monitor"
	"github.com/alexpfx/go_process_monitor/util"
	"github.com/urfave/cli/v2"
	"strings"

	"os"
	"regexp"
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

func main() {
	app := cli.App{

		Commands: []*cli.Command{

			{
				Name: "listen",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "porta onde deve esperar conexões http",
						Value:   7777,
					},
					&cli.StringFlag{
						Name:     "ps",
						Usage:    "processo que terá sua saída monitorada",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "auth",
						Aliases:  []string{"a"},
						Usage:    "user:pass",
						Value: "test:test00",
						Required: true,
					},
				},

				Action: func(c *cli.Context) error {
					auth := strings.Split(c.String("auth"), ":")
					if len(auth) != 2 {
						return fmt.Errorf("user e pass são obrigatórios")
					}

					psSplitted := strings.Split(c.String("ps"), " ")
					cmdPath := psSplitted[0]
					cmdArgs := psSplitted[1:]

					queue := make(map[string]bool)
					ps := monitor.Process{
						CmdPath: cmdPath,
						Args:    cmdArgs,
					}

					psCh := ps.Start()

					srv := monitor.Server{
						Port: c.Int("port"),
						User: auth[0],
						Pass: auth[1],
					}

					srvCh := srv.Start()

					for {
						select {
						case psMsg := <-psCh:
							which := search(queue, psMsg)
							if which != "" {
								queue[which] = false
								fmt.Printf(NoticeColor+"\n", which)
							} else {
								//fmt.Println(psMsg)
							}
						case srMsg := <-srvCh:
							queue[srMsg] = !queue[srMsg]
							fmt.Printf("%v ", queue)
						}
					}
				},
			},
		},
	}

	err := app.Run(os.Args)
	util.Check(err)

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

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
