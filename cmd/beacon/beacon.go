package main

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/commands"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	app := cli.NewCLI("beacon", "1.0.0")
	// init cli folder and templates
	startup.StartCli()
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"audit": func() (cli.Command, error) {
			return &commands.K8sBenchmark{}, nil
		},
		"a": func() (cli.Command, error) {
			return &commands.K8sBenchmark{}, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}
