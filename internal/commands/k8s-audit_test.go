package commands

import (
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/mocks"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

//Test_EvalVarSingleIn text
func Test_EvalVarSingleIn(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiProcessInClause.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarSingleNotInGood text
func Test_EvalVarSingleNotInGood(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiProcessInClause.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"ttt,aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarSingleNotInBad text
func Test_EvalVarSingleNotInBad(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiProcessInClause.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"RBAC,aaa"}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarSingleNotInSingleValue test
func Test_EvalVarSingleNotInSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiProcessInClause.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiExprProcessParam.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"AlwaysAdmit"}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprMultiValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiExprProcessParam.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"bbb,aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarMultiExprMultiEmptyValue test
func Test_EvalVarMultiExprMultiEmptyValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeMultiExprEmptyProcessParam.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{common.GrepRegex}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarComparator test
func Test_EvalVarComparator(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckTypeComparator.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"1204"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_MultiCommandParams_OK test
func Test_MultiCommandParams_OK(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamOK.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParams_OK_With_IN test
func Test_MultiCommandParams_OK_With_IN(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamOKWithIN.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk,aaa"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParams_NOK test
func Test_MultiCommandParams_NOK(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamNOK.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParams_NOKWith_IN test
func Test_MultiCommandParams_NOKWith_IN(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamNOKWithIN.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParamsPass1stResultToNext test
func Test_MultiCommandParamsPass1stResultToNext(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamPass1stResultToNext.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("ccc kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParamsComplex test
func Test_MultiCommandParamsComplex(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckMultiParamComplex.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "/etc/kubernetes/pki/encry.yaml"}, nil).Times(1)
	executor.EXPECT().Exec("bbb").Return(&shell.CommandResult{Stdout: "/etc/kubernetes/pki/encry.yaml"}, nil).Times(1)
	executor.EXPECT().Exec("ccc").Return(&shell.CommandResult{Stdout: "aescbc"}, nil).Times(1)
	executor.EXPECT().Exec("ddd").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	executor.EXPECT().Exec("eee").Return(&shell.CommandResult{Stdout: "secretbox"}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParamsComplexOppositeEmptyReturn test
func Test_MultiCommandParamsComplexOppositeEmptyReturn(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckInClauseOppositeEmptyReturn.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommandParamsComplexOppositeWithNumber test
func Test_MultiCommandParamsComplexOppositeWithNumber(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckInClauseOppositeWithNum.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_MultiCommand4_2_13 test
func Test_MultiCommand4_2_13(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal(readTestData("CheckInClause4.2.13.json", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("ps -ef | grep kubelet |grep ' --config' | grep -o ' --config=[^\"]\\S*' | awk -F \"=\" '{print $2}' |awk 'FNR <= 1'").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	executor.EXPECT().Exec("ps -ef | grep kubelet |grep 'TLSCipherSuites' | grep -o 'TLSCipherSuites=[^\"]\\S*' | awk -F \"=\" '{print $2}' |awk 'FNR <= 1'").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	kb := K8sAudit{Command: executor, resultProcessor: getResultProcessingFunction([]string{})}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
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
