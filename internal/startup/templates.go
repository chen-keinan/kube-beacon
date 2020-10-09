package startup

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/gobuffalo/packr"
	"os"
	"path/filepath"
)

//LoadK8sBenchmarkFile use packr to load benchmark audit test json
func LoadK8sBenchmarkFile() []BenchFilesData {
	box := packr.NewBox("./../benchmark/k8s")
	s, err := box.FindString(common.MasterNodeConfigurationFiles)
	a, err := box.FindString(common.APIServer)
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	return []BenchFilesData{{common.MasterNodeConfigurationFiles, s},
		{common.APIServer, a}}
}

//CreateBenchmarkFileIfNotExist create benchmark audit file if not exist
func CreateBenchmarkFileIfNotExist(filesData []BenchFilesData) error {
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
				fmt.Printf("faild to close file %s", filePath)
			}
		}
	}
	return nil
}

type BenchFilesData struct {
	Name string
	Data string
}
