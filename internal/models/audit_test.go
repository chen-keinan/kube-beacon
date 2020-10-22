package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

//Test_CheckType_Permission test
func Test_CheckType_Blaa(t *testing.T) {
	ab := Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeBlaa.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Nil(t, ab.Categories[0].SubCategory.AuditTests[0].CmdExprBuilder)
}

//Test_CheckType_Multi_ProcessParam test
func Test_CheckType_Multi_ProcessParam(t *testing.T) {
	ab := Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiProcessParam.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, ab.Categories[0].SubCategory.AuditTests[0].CmdExprBuilder)
}

func readTestData(fileName string, t *testing.T) []byte {
	f, err := os.Open(fmt.Sprintf("./fixtures/%s", fileName))
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return b
}

//Test_AuditIndexWithParams test
func Test_AuditIndexWithParams(t *testing.T) {
	commands := []string{"aaa", "bbb #0", "ccc #0 #1"}
	commandParams := make(map[int][]string)
	for index, command := range commands {
		findIndex(command, "#", index, commandParams)
	}
	assert.Equal(t, commandParams[1], []string{"0"})
	assert.Equal(t, commandParams[2], []string{"0", "1"})
}

//Test_AuditIndexWithParams test
func Test_AuditIndexWithNoParams(t *testing.T) {
	commands := []string{"aaa", "bbb", "ccc"}
	commandParams := make(map[int][]string)
	for index, command := range commands {
		findIndex(command, "#", index, commandParams)
	}
	assert.True(t, len(commandParams) == 0)
}
