package cli

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//Test_StartCli tests
func Test_StartCli(t *testing.T) {
	StartCli()
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
	sArgs := ArgsSanitizer(args)
	assert.Equal(t, sArgs[0], "a")
	assert.Equal(t, sArgs[1], "b")
	args = []string{}
	sArgs = ArgsSanitizer(args)
	assert.True(t, sArgs[0] == "")
}

//Test_BeaconHelpFunc test
func Test_BeaconHelpFunc(t *testing.T) {
	cm := make(map[string]cli.CommandFactory)
	bhf := BeaconHelpFunc("Beacon")
	helpFile := bhf(cm)
	assert.True(t, strings.Contains(helpFile, "Available commands are:"))
	assert.True(t, strings.Contains(helpFile, "Usage: Beacon [--version] [--help] <command> [<args>]"))
}
