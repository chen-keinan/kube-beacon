package startup

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
)

//GenerateK8sBenchmarkFiles use packr to load benchmark audit test yaml
func GenerateK8sBenchmarkFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../benchmark/k8s/v1.6.0/")
	// Add Master Node Configuration tests
	mnc, err := box.FindString(common.MasterNodeConfigurationFiles)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.MasterNodeConfigurationFiles, Data: mnc})
	// Add API Server tests
	aps, err := box.FindString(common.APIServer)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.APIServer, Data: aps})
	// Add Controller Manager tests
	cm, err := box.FindString(common.ControllerManager)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ControllerManager, Data: cm})
	// Add Scheduler tests
	sc, err := box.FindString(common.Scheduler)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Scheduler, Data: sc})
	// Add Control Plane Configuration tests
	etcd, err := box.FindString(common.Etcd)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Etcd, Data: etcd})
	// Add Etcd tests
	cpc, err := box.FindString(common.ControlPlaneConfiguration)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.ControlPlaneConfiguration, Data: cpc})
	// Add Worker Nodes tests
	wn, err := box.FindString(common.WorkerNodes)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.WorkerNodes, Data: wn})
	p, err := box.FindString(common.Policies)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.Policies, Data: p})
	return fileInfo, nil
}

//GenerateGkeBenchmarkFiles use packr to load benchmark audit test yaml
func GenerateGkeBenchmarkFiles() ([]utils.FilesInfo, error) {
	fileInfo := make([]utils.FilesInfo, 0)
	box := packr.NewBox("./../benchmark/gke/v1.1.0/")
	// Add Master Node Configuration tests
	cpc, err := box.FindString(common.GkeControlPlaneConfiguration)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.GkeControlPlaneConfiguration, Data: cpc})
	// Add Worker Nodes tests
	wn, err := box.FindString(common.GkeWorkerNodes)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.GkeWorkerNodes, Data: wn})
	gkp, err := box.FindString(common.GkePolicies)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.GkePolicies, Data: gkp})
	gms, err := box.FindString(common.GkeManagedServices)
	if err != nil {
		return []utils.FilesInfo{}, fmt.Errorf("faild to load k8s benchmarks audit tests %s", err.Error())
	}
	fileInfo = append(fileInfo, utils.FilesInfo{Name: common.GkeManagedServices, Data: gms})
	return fileInfo, nil
}

//GetHelpSynopsis get help synopsis file
func GetHelpSynopsis() string {
	box := packr.NewBox("./../cli/commands/help/")
	// Add Master Node Configuration tests
	hs, err := box.FindString(common.Synopsis)
	if err != nil {
		panic(fmt.Sprintf("faild to load cli help synopsis %s", err.Error()))
	}
	return hs
}

//SaveBenchmarkFilesIfNotExist create benchmark audit file if not exist
func SaveBenchmarkFilesIfNotExist(spec, version string, filesData []utils.FilesInfo) error {
	for _, fileData := range filesData {
		filePath := filepath.Join(utils.GetBenchmarkFolder(spec, version), fileData.Name)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf(err.Error())
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
