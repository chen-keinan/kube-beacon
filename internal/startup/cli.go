package startup

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
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
	benchK8s := LoadK8sBenchmarkFile()
	err = utils.CreateBenchmarkFileIfNotExist(common.K8sBenchmarkAuditFileName, benchK8s)
	if err != nil {
		panic(err)
	}
}
