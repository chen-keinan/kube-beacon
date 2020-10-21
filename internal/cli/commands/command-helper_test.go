package commands

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

//Test_AddFailedMessages text
func Test_AddFailedMessages(t *testing.T) {
	atb1 := &models.AuditBench{TestResult: &models.AuditResult{NumOfExec: 1, NumOfSuccess: 0}}
	ve1 := ValidateExprData{atb: atb1}
	afm := AddFailedMessages(ve1)
	assert.True(t, len(afm) == 1)
	atbr2 := &models.AuditBench{TestResult: &models.AuditResult{NumOfExec: 1, NumOfSuccess: 1}}
	ve2 := ValidateExprData{atb: atbr2}
	afm = AddFailedMessages(ve2)
	assert.True(t, len(afm) == 0)
}

//Test_isArgsExist
func Test_isArgsExist(t *testing.T) {
	args := []string{"aaa", "bbb"}
	exist := isArgsExist(args, "aaa")
	assert.True(t, exist)
	exist = isArgsExist(args, "ccc")
	assert.False(t, exist)
}

//Test_isArgsExist
func Test_GetProcessingFunction(t *testing.T) {
	args := []string{"r"}
	a := getResultProcessingFunction(args)
	name := GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func2"))
	args = []string{}
	a = getResultProcessingFunction(args)
	name = GetFunctionName(a)
	assert.True(t, strings.Contains(name, "commands.glob..func1"))
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
