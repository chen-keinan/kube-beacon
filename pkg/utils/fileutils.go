package utils

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

//GetHomeFolder return beacon home folder
func GetHomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic("Failed to fetch user home folder")
	}
	return path.Join(usr.HomeDir, ".beacon")
}

//CreateHomeFolderIfNotExist create beacon home folder if not exist
func CreateHomeFolderIfNotExist() error {
	beaconFolder := GetHomeFolder()
	_, err := os.Stat(beaconFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(beaconFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon home folder at %s", beaconFolder)
		}
	}
	return nil
}

//GetBenchmarkFolder return benchmark folder
func GetBenchmarkFolder() string {
	return filepath.Join(GetHomeFolder(), "benchmarks")
}

//CreateBenchmarkFolderIfNotExist create beacon benchmark folder if not exist
func CreateBenchmarkFolderIfNotExist() error {
	benchmarkFolder := GetBenchmarkFolder()
	_, err := os.Stat(benchmarkFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(benchmarkFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon benchmark folder folder at %s", benchmarkFolder)
		}
	}
	return nil
}

//CreateBenchmarkFileIfNotExist create benchmark audit file if not exist
func CreateBenchmarkFileIfNotExist(filename, fileData string) error {
	filePath := filepath.Join(GetBenchmarkFolder(), filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(fileData)
		if err != nil {
			return fmt.Errorf("failed to write benchmark file")
		}
		err = f.Close()
		if err != nil {
			fmt.Printf("faild to close file %s", filePath)
		}
	}
	return nil
}

//GetK8sBenchmarkAuditTestsFile return k8s benchmark file
func GetK8sBenchmarkAuditTestsFile() []string {
	filePath := filepath.Join(GetBenchmarkFolder(), filepath.Clean(common.MasterNodeConfigurationFiles))
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		panic("failed to read k8s benchmark audit file")
	}
	return []string{string(data)}
}
