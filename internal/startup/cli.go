package startup

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/commands"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

//StartCli init beacon cli , folder , templates and etc
func StartCli() {
	err := utils.CreateHomeFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	filesData := LoadK8sBenchmarkFile()
	err = CreateBenchmarkFileIfNotExist(filesData)
	if err != nil {
		panic(err)
	}
}

//InitCLI initialize beacon cli
func InitCLI(sa SanitizeArgs) {
	app := cli.NewCLI("beacon", "1.0.0")
	// init cli folder and templates
	StartCli()
	app.Args = sa(os.Args[1:])
	app.Commands = map[string]cli.CommandFactory{
		"audit": func() (cli.Command, error) {
			return &commands.K8sAudit{}, nil
		},
		"a": func() (cli.Command, error) {
			return &commands.K8sAudit{}, nil
		},
	}
	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}

//ArgsSanitizer sanitize CLI arguments
var ArgsSanitizer SanitizeArgs = func(str []string) []string {
	args := make([]string, 0)
	for _, arg := range str {
		arg = strings.Replace(arg, "--", "", -1)
		arg = strings.Replace(arg, "-", "", -1)
		args = append(args, arg)
	}
	return args
}

//SanitizeArgs sanitizer func
type SanitizeArgs func(str []string) []string
