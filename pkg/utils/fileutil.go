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
	// User can set a custom BEACON_HOME from environment variable
	usrHome := GetEnv(common.BeaconHomeEnvVar, usr.HomeDir)
	return path.Join(usrHome, ".beacon")
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
func GetBenchmarkFolder(spec, version string) string {
	return filepath.Join(GetHomeFolder(), fmt.Sprintf("benchmarks/%s/%s/", spec, version))
}

//CreateBenchmarkFolderIfNotExist create beacon benchmark folder if not exist
func CreateBenchmarkFolderIfNotExist(spec, version string) error {
	benchmarkFolder := GetBenchmarkFolder(spec, version)
	_, err := os.Stat(benchmarkFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(benchmarkFolder, 0750)
		if errDir != nil {
			return fmt.Errorf("failed to create beacon benchmark folder folder at %s", benchmarkFolder)
		}
	}
	return nil
}

//GetK8sBenchAuditFiles return k8s benchmark file
func GetK8sBenchAuditFiles(spec, version string) ([]FilesInfo, error) {
	filesData := make([]FilesInfo, 0)
	folder := GetBenchmarkFolder(spec, version)
	filesInfo, err := ioutil.ReadDir(filepath.Join(folder))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range filesInfo {
		filePath := filepath.Join(GetBenchmarkFolder(spec, version), filepath.Clean(fileInfo.Name()))
		fData, err := ioutil.ReadFile(filepath.Clean(filePath))
		if err != nil {
			return nil, err
		}
		filesData = append(filesData, FilesInfo{fileInfo.Name(), string(fData)})
	}
	return filesData, nil
}

//FilesInfo file data
type FilesInfo struct {
	Name string
	Data string
}

//GetEnv Get Environment Variable value or return default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
