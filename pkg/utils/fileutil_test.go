package utils

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Test_GetHomeFolder test
func Test_GetHomeFolder(t *testing.T) {
	a := GetHomeFolder()
	assert.True(t, strings.HasSuffix(a, ".beacon"))
}

//Test_CreateHomeFolderIfNotExist test
func Test_CreateHomeFolderIfNotExist(t *testing.T) {
	err := CreateHomeFolderIfNotExist()
	assert.NoError(t, err)
	_, err = os.Stat(GetHomeFolder())
	if os.IsNotExist(err) {
		t.Fatal()
	}
	err = os.RemoveAll(GetHomeFolder())
	assert.NoError(t, err)
}

//Test_GetBenchmarkFolder test
func Test_GetBenchmarkFolder(t *testing.T) {
	a := GetBenchmarkFolder()
	assert.True(t, strings.HasSuffix(a, ".beacon/benchmarks/v1.6.0"))
}

//Test_CreateBenchmarkFolderIfNotExist test
func Test_CreateBenchmarkFolderIfNotExist(t *testing.T) {
	err := CreateBenchmarkFolderIfNotExist()
	assert.NoError(t, err)
	_, err = os.Stat(GetBenchmarkFolder())
	if os.IsNotExist(err) {
		t.Fatal()
	}
	err = os.RemoveAll(GetBenchmarkFolder())
	assert.NoError(t, err)
}

//Test_GetK8sBenchAuditFiles test
func Test_GetK8sBenchAuditFiles(t *testing.T) {
	err := CreateHomeFolderIfNotExist()
	if err != nil {
		t.Fatal(err)
	}
	err = CreateBenchmarkFolderIfNotExist()
	if err != nil {
		t.Fatal(err)
	}
	err = saveFilesIfNotExist([]FilesInfo{{Name: "aaa", Data: "bbb"}, {Name: "ddd", Data: "ccc"}})
	if err != nil {
		t.Fatal(err)
	}
	f, err := GetK8sBenchAuditFiles()
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll(GetBenchmarkFolder())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, f[0].Name, "aaa")
	assert.Equal(t, f[1].Name, "ddd")

}

//Test_GetK8sBenchAuditNoFolder test
func Test_GetK8sBenchAuditNoFolder(t *testing.T) {
	_, err := GetK8sBenchAuditFiles()
	assert.True(t, err != nil)
}

func saveFilesIfNotExist(filesData []FilesInfo) error {
	for _, fileData := range filesData {
		filePath := filepath.Join(GetBenchmarkFolder(), fileData.Name)
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

//Test_GetEnv test getting home beacon folder
func Test_GetEnv(t *testing.T) {
	os.Setenv(common.BeaconHomeEnvVar, "/home/beacon")
	homeFolder := GetEnv(common.BeaconHomeEnvVar, "/home/user")
	assert.Equal(t, homeFolder, "/home/beacon")
	os.Unsetenv(common.BeaconHomeEnvVar)
	homeFolder = GetEnv(common.BeaconHomeEnvVar, "/home/user")
	assert.Equal(t, homeFolder, "/home/user")
}
