package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
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
