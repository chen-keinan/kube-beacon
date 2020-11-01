package cli

import (
	"github.com/chen-keinan/beacon/internal/cli/commands"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/mocks"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

//Test_StartCli tests
func Test_StartCli(t *testing.T) {
	InitCli()
	files, err := utils.GetK8sBenchAuditFiles()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(files), 8)
	assert.Equal(t, files[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, files[1].Name, common.APIServer)
	assert.Equal(t, files[2].Name, common.ControllerManager)
	assert.Equal(t, files[3].Name, common.Scheduler)
	assert.Equal(t, files[4].Name, common.Etcd)
	assert.Equal(t, files[5].Name, common.ControlPlaneConfiguration)
	assert.Equal(t, files[6].Name, common.WorkerNodes)
	assert.Equal(t, files[7].Name, common.Policies)
}

func Test_ArgsSanitizer(t *testing.T) {
	args := []string{"--a", "-b"}
	sArgs, helpNeed := ArgsSanitizer(args)
	assert.Equal(t, sArgs[0], "a")
	assert.Equal(t, sArgs[1], "b")
	assert.False(t, helpNeed)
	args = []string{}
	sArgs, _ = ArgsSanitizer(args)
	assert.True(t, sArgs[0] == "")
	args = []string{"--help"}
	_, helpNeed2nd := ArgsSanitizer(args)
	assert.True(t, helpNeed2nd)
}

//Test_BeaconHelpFunc test
func Test_BeaconHelpFunc(t *testing.T) {
	cm := make(map[string]cli.CommandFactory)
	bhf := BeaconHelpFunc("Beacon")
	helpFile := bhf(cm)
	assert.True(t, strings.Contains(helpFile, "Available commands are:"))
	assert.True(t, strings.Contains(helpFile, "Usage: Beacon [--version] [--help] <command> [<args>]"))
}

//Test_createCliBuilderData test
func Test_createCliBuilderData(t *testing.T) {
	cmdArgs := []string{"a"}
	cliArgs, _ := ArgsSanitizer(os.Args[1:])
	cmdArgs = append(cmdArgs, cliArgs...)
	cmds := make([]cli.Command, 0)
	// invoke cli
	cmds = append(cmds, commands.NewK8sAudit(cmdArgs))
	c := createCliBuilderData(cmdArgs, cmds)
	_, ok := c["a"]
	assert.True(t, ok)

}

//Test_InvokeCli test
func Test_InvokeCli(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa"}
	ab.EvalExpr = "'$0' != '';"
	ab.CommandParams = map[int][]string{}
	ab.CmdExprBuilder = utils.UpdateCmdExprParam
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "1234"}, nil).Times(1)
	tl := mocks.NewMockTestLoader(ctrl)
	tl.EXPECT().LoadAuditTests().Return([]*models.SubCategory{{Name: "te", AuditTests: []*models.AuditBench{ab}}})
	kb := &commands.K8sAudit{Command: executor, ResultProcessor: commands.GetResultProcessingFunction([]string{}), FileLoader: tl, OutputGenerator: commands.ConsoleOutputGenerator}
	cmdArgs := []string{"a"}
	cmds := make([]cli.Command, 0)
	// invoke cli
	cmds = append(cmds, kb)
	c := createCliBuilderData(cmdArgs, cmds)
	a, err := invokeCommandCli(cmdArgs, c)
	assert.NoError(t, err)
	assert.True(t, a == 0)
}
