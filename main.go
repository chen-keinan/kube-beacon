package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	app := cli.NewCLI("beacon", "1.0.0")
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"test": func() (cli.Command, error) {
			return &Hello{}, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}

type Hello struct {
}

func (*Hello) Help() string {
	return "-t , --test run benchmark tests"
}
func (*Hello) Run(args []string) int {
	fmt.Printf("running benchmark tests, %v", args)
	return 0
}
func (h *Hello) Synopsis() string {
	return h.Help()
}
