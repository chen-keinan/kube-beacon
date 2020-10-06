package beacon

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("beacon", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{

	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
