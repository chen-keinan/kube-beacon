package models

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

//Test_CheckType_Permission test
func Test_CheckType_Blaa(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal(readTestData("CheckTypeBlaa.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Nil(t, ab.Categories[0].SubCategory.AuditTests[0].Sanitize)
}

//Test_CheckType_Multi_ProcessParam test
func Test_CheckType_Multi_ProcessParam(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiProcessParam.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, ab.Categories[0].SubCategory.AuditTests[0].Sanitize)
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
