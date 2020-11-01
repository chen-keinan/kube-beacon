package cli

import (
	"bytes"
	"fmt"
	"github.com/chen-keinan/beacon/internal/cli/commands"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

var log = logger.GetLog()

//InitCli init beacon cli , folder , templates and etc
func InitCli() {
	err := utils.CreateHomeFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	filesData, err := startup.GenerateK8sBenchmarkFiles()
	if err != nil {
		panic(err)
	}
	err = startup.SaveBenchmarkFilesIfNotExist(filesData)
	if err != nil {
		panic(err)
	}
}

//StartCLI initialize beacon cli
func StartCLI(sa SanitizeArgs) {
	// create cli data
	cmdArgs := []string{"a"}
	cliArgs, helpNeeded := sa(os.Args[1:])
	cmds := make([]cli.Command, 0)
	cmdArgs = append(cmdArgs, cliArgs...)
	// invoke cli
	cmds = append(cmds, commands.NewK8sAudit(cliArgs))
	commands := createCliBuilderData(cmdArgs, cmds)
	if helpNeeded {
		cmdArgs = cmdArgs[1:]
	}
	status, err := invokeCommandCli(cmdArgs, commands)
	if err != nil {
		log.Console(err.Error())
	}
	os.Exit(status)
}

//createCliBuilderData return cli params and commands
func createCliBuilderData(ca []string, cmd []cli.Command) map[string]cli.CommandFactory {
	// read cli args
	cmdFactory := make(map[string]cli.CommandFactory)
	// build cli commands
	for index, a := range cmd {
		cmdFactory[ca[index]] = func() (cli.Command, error) {
			return a, nil
		}
	}
	return cmdFactory
}

// invokeCommandCli invoke cli command with params
func invokeCommandCli(args []string, commands map[string]cli.CommandFactory) (int, error) {
	app := cli.NewCLI(common.BeaconCli, common.BeaconVersion)
	// init cli folder and templates
	InitCli()
	app.Args = append(app.Args, args...)
	app.Commands = commands
	app.HelpFunc = BeaconHelpFunc(common.BeaconCli)
	status, err := app.Run()
	return status, err
}

//ArgsSanitizer sanitize CLI arguments
var ArgsSanitizer SanitizeArgs = func(str []string) ([]string, bool) {
	var helpNeeded bool
	args := make([]string, 0)
	if len(str) == 0 {
		args = append(args, "")
	}
	for _, arg := range str {
		arg = strings.Replace(arg, "--", "", -1)
		arg = strings.Replace(arg, "-", "", -1)
		args = append(args, arg)
		if arg == "help" {
			helpNeeded = true
		}
	}
	return args, helpNeeded
}

//SanitizeArgs sanitizer func
type SanitizeArgs func(str []string) ([]string, bool)

// BeaconHelpFunc beacon help function with all supported commands
func BeaconHelpFunc(app string) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(startup.GetHelpSynopsis(), app))
		return buf.String()
	}
}
