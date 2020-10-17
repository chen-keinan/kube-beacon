package startup

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
)

//GenerateK8sBenchmarkFiles use packr to load benchmark audit test json
func GenerateK8sBenchmarkFiles() []utils.FilesInfo {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../benchmark/k8s/v1.6.0/")
	// Add Master Node Configuration tests
	mnc, err := box.FindString(common.MasterNodeConfigurationFiles)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.MasterNodeConfigurationFiles, Data: mnc})
	// Add API Server tests
	aps, err := box.FindString(common.APIServer)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.APIServer, Data: aps})
	// Add Controller Manager tests
	cm, err := box.FindString(common.ControllerManager)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ControllerManager, Data: cm})
	// Add Scheduler tests
	sc, err := box.FindString(common.Scheduler)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Scheduler, Data: sc})
	// Add Etcd tests
	etcd, err := box.FindString(common.Etcd)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Etcd, Data: etcd})
	return fileInfo
}

//SaveBenchmarkFilesIfNotExist create benchmark audit file if not exist
func SaveBenchmarkFilesIfNotExist(filesData []utils.FilesInfo) error {
	for _, fileData := range filesData {
		filePath := filepath.Join(utils.GetBenchmarkFolder(), fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			_, err = f.WriteString(fileData.Data)
			if err != nil {
				return fmt.Errorf("failed to write benchmark file")
			}
			err = f.Close()
			if err != nil {
				return fmt.Errorf("faild to close file %s", filePath)
			}
		}
	}
	return nil
}
