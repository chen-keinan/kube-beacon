package startup

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

//Test_CreateBenchmarkFilesIfNotExist test
func Test_CreateBenchmarkFilesIfNotExist(t *testing.T) {
	bFiles := GenerateK8sBenchmarkFiles()
	// generate test with packr
	assert.Equal(t, bFiles[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, bFiles[1].Name, common.APIServer)
	assert.Equal(t, bFiles[2].Name, common.ControllerManager)
	err := utils.CreateBenchmarkFolderIfNotExist()
	assert.NoError(t, err)
	// save benchmark files to folder
	err = SaveBenchmarkFilesIfNotExist(bFiles)
	assert.NoError(t, err)
	// fetch files from benchmark folder
	bFiles, err = utils.GetK8sBenchAuditFiles()
	assert.Equal(t, bFiles[0].Name, common.MasterNodeConfigurationFiles)
	assert.Equal(t, bFiles[1].Name, common.APIServer)
	assert.Equal(t, bFiles[2].Name, common.ControllerManager)
	assert.NoError(t, err)
	err = os.RemoveAll(utils.GetBenchmarkFolder())
	assert.NoError(t, err)
}
