package main

import (
	"fmt"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/kyokomi/emoji"
	"github.com/mitchellh/cli"
	"os"
	"strconv"
	"strings"
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
	ls := execute.ExecTask{
		Command: "stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml",
		Args:    []string{},
		Shell:   true,
	}
	res, err := ls.Execute()
	if err != nil {
		fmt.Sprintf("Failed to execute command %s", err.Error())
	}
	value,err:= strconv.Atoi(strings.Replace(res.Stderr,"stdout: ","",-1))
	if err != nil{
		fmt.Println(res.Stderr)
		//fmt.Print("failed to convert string %s",err.Error())
	}
	if value <= 644 {
		fmt.Print("ok\n")
		emoji.Println(":OK_hand:Test pass!!!")
	}
 	return 0
}
func (h *Hello) Synopsis() string {
	return h.Help()
}
