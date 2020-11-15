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

//InitBenchmarkSpecData initialize benchmark spec file and save if to file system
func InitBenchmarkSpecData(spec, version string) {
	err = utils.CreateHomeFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist(spec, version)
	if err != nil {
		panic(err)
	}
	var filesData []utils.FilesInfo
	switch spec {
	case "k8s":
		if version == "v1.6.0" {
			filesData, err = startup.GenerateK8sBenchmarkFiles()
		}
	case "gks":
		if version == "v1.1.0" {
			filesData, err = startup.GenerateGksBenchmarkFiles()
		}
	}
	if err != nil {
		panic(err)
	}
	err = startup.SaveBenchmarkFilesIfNotExist(spec, version, filesData)
	if err != nil {
		panic(err)
	}
}

//InitBenchmarkSpecData init beacon cli , folder , templates and etc
var err error

//StartCLI initialize beacon cli
func StartCLI(sa SanitizeArgs) {
	// create cli data
	cmdArgs := []string{"a"}
	ad := sa(os.Args[1:])
	cmds := make([]cli.Command, 0)
	cmdArgs = append(cmdArgs, ad.filters...)
	// invoke cli
	cmds = append(cmds, commands.NewK8sAudit(ad.filters, ad.specType, ad.specVersion))
	commands := createCliBuilderData(cmdArgs, cmds)
	if ad.help {
		cmdArgs = cmdArgs[1:]
	}
	// init cli folder and templates
	InitBenchmarkSpecData(ad.specType, ad.specVersion)
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
	app.Args = append(app.Args, args...)
	app.Commands = commands
	app.HelpFunc = BeaconHelpFunc(common.BeaconCli)
	status, err := app.Run()
	return status, err
}

//ArgsSanitizer sanitize CLI arguments
var ArgsSanitizer SanitizeArgs = func(str []string) ArgsData {
	ad := ArgsData{specType: "k8s", specVersion: "v1.6.0"}
	args := make([]string, 0)
	if len(str) == 0 {
		args = append(args, "")
	}
	for _, arg := range str {
		arg = strings.Replace(arg, "--", "", -1)
		arg = strings.Replace(arg, "-", "", -1)
		switch {
		case arg == "help", arg == "h":
			ad.help = true
			args = append(args, arg)
		case strings.HasPrefix(arg, "s="):
			ad.specType = arg[len("s="):]
		case strings.HasPrefix(arg, "v="):
			ad.specVersion = fmt.Sprintf("v%s", arg[len("v="):])
		default:
			args = append(args, arg)
		}
	}
	ad.filters = args
	return ad
}

//ArgsData hold cli args data
type ArgsData struct {
	filters     []string
	help        bool
	specType    string
	specVersion string
}

//SanitizeArgs sanitizer func
type SanitizeArgs func(str []string) ArgsData

// BeaconHelpFunc beacon help function with all supported commands
func BeaconHelpFunc(app string) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(startup.GetHelpSynopsis(), app))
		return buf.String()
	}
}
