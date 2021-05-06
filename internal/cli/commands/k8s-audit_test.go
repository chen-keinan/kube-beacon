package commands

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/mocks"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/shell"
	m2 "github.com/chen-keinan/beacon/pkg/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

//Test_EvalVarSingleIn text
func Test_EvalVarSingleIn(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiProcessInClause.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"aaa"}, 1, make([]string, 0), 0)
	assert.True(t, k == 0)
}

//Test_EvalVarSingleNotInGood text
func Test_EvalVarSingleNotInGood(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiProcessInClause.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"ttt,aaa"}, 1, make([]string, 0), 0)
	assert.True(t, k == 0)
}

//Test_EvalVarSingleNotInBad text
func Test_EvalVarSingleNotInBad(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiProcessInClause.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"RBAC,aaa"}, 1, make([]string, 0), 0)
	assert.True(t, k > 0)
}

//Test_EvalVarSingleNotInSingleValue test
func Test_EvalVarSingleNotInSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiProcessInClause.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"aaa"}, 1, make([]string, 0), 0)
	assert.True(t, k == 0)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiExprProcessParam.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"AlwaysAdmit"}, 1, make([]string, 0), 0)
	assert.True(t, k > 0)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprMultiValue(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiExprProcessParam.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"bbb,aaa"}, 1, make([]string, 0), 0)
	assert.True(t, k == 0)
}

//Test_EvalVarMultiExprMultiEmptyValue test
func Test_EvalVarMultiExprMultiEmptyValue(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeMultiExprEmptyProcessParam.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{common.GrepRegex}, 1, make([]string, 0), 0)
	assert.True(t, k > 0)
}

//Test_EvalVarComparator test
func Test_EvalVarComparator(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckTypeComparator.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	k := kb.evalExpression(bench, []string{"1204"}, 1, make([]string, 0), 0)
	assert.True(t, k == 0)
}

//Test_MultiCommandParams_OK test
func Test_MultiCommandParams_OK(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamOK.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommandParams_OK_With_IN test
func Test_MultiCommandParams_OK_With_IN(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamOKWithIN.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk,aaa"}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommandParams_NOK test
func Test_MultiCommandParams_NOK(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamNOK.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommandParams_NOKWith_IN test
func Test_MultiCommandParams_NOKWith_IN(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamNOKWithIN.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{})}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
}

//Test_MultiCommandParamsPass1stResultToNext test
func Test_MultiCommandParamsPass1stResultToNext(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamPass1stResultToNext.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("ccc kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{})}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
}

//Test_MultiCommandParamsComplex test
func Test_MultiCommandParamsComplex(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckMultiParamComplex.yml", t), &ab)
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
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{})}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
}

//Test_MultiCommandParamsComplexOppositeEmptyReturn test
func Test_MultiCommandParamsComplexOppositeEmptyReturn(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckInClauseOppositeEmptyReturn.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommandParamsComplexOppositeWithNumber test
func Test_MultiCommandParamsComplexOppositeWithNumber(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckInClauseOppositeWithNum.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommand4_2_13 test
func Test_MultiCommand4_2_13(t *testing.T) {
	ab := models.Audit{}
	err := yaml.Unmarshal(readTestData("CheckInClause4.2.13.yml", t), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("ps -ef | grep kubelet |grep ' --config' | grep -o ' --config=[^\"]\\S*' | awk -F \"=\" '{print $2}' |awk 'FNR <= 1'").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	executor.EXPECT().Exec("ps -ef | grep kubelet |grep 'TLSCipherSuites' | grep -o 'TLSCipherSuites=[^\"]\\S*' | awk -F \"=\" '{print $2}' |awk 'FNR <= 1'").Return(&shell.CommandResult{Stdout: ""}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_MultiCommand4_2_13 test
func Test_MultiCommand5_2_7(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa", "bbb #0"}
	ab.EvalExpr = "'NET_RAW' IN ($1); || 'ALL' IN ($1);"
	ab.CommandParams = map[int][]string{1: {"0"}}
	ab.CmdExprBuilder = utils.UpdateCmdExprParam
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "1234\n"}, nil).Times(1)
	executor.EXPECT().Exec("bbb 1234").Return(&shell.CommandResult{Stdout: "\n"}, nil).Times(1)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab)
	assert.False(t, ab.TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()

}

//Test_MultiCommand5_3_4 test
func Test_MultiCommand5_3_4(t *testing.T) {
	ab := &models.AuditBench{}
	ab.AuditCommand = []string{"aaa", "bbb"}
	ab.EvalExpr = "'$0' == ''; && '$1' == '';"
	ab.CommandParams = map[int][]string{}
	ab.CmdExprBuilder = utils.UpdateCmdExprParam
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "\n\n\n\n\n"}, nil).Times(1)
	executor.EXPECT().Exec("bbb").Return(&shell.CommandResult{Stdout: "default-token-ppzx7\n\n\n\n\n"}, nil).Times(1)
	kb := K8sAudit{Command: executor, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: plChan, CompletedChan: completedChan}
	kb.runAuditTest(ab)
	assert.False(t, ab.TestSucceed)
	go func() {
		<-plChan
		completedChan <- true
	}()
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

//Test_NewK8sAudit test
func Test_NewK8sAudit(t *testing.T) {
	args := []string{"a", "i=1.2.3"}
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	ka := NewK8sAudit(args, "k8s", "v1.6.0", plChan, completedChan)
	assert.True(t, len(ka.PredicateParams) == 2)
	assert.True(t, len(ka.PredicateChain) == 2)
	assert.True(t, ka.ResultProcessor != nil)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_Help test
func Test_Help(t *testing.T) {
	args := []string{"a", "i=1.2.3"}
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	ka := NewK8sAudit(args, "k8s", "v1.6.0", plChan, completedChan)
	help := ka.Help()
	assert.True(t, len(help) > 0)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

//Test_reportResultProcessor test
func Test_reportResultProcessor(t *testing.T) {
	ad := &models.AuditBench{Name: "1.2.1 aaa"}
	fm := reportResultProcessor(ad, 0)
	assert.True(t, len(fm) == 0)
	fm = reportResultProcessor(ad, 1)
	assert.True(t, len(fm) == 1)
	assert.Equal(t, fm[0].Name, "1.2.1 aaa")
}

//Test_K8sSynopsis test
func Test_K8sSynopsis(t *testing.T) {
	args := []string{"a", "i=1.2.3"}
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	ka := NewK8sAudit(args, "k8s", "v1.6.0", plChan, completedChan)
	s := ka.Synopsis()
	assert.True(t, len(s) > 0)
	go func() {
		<-plChan
		completedChan <- true
	}()
}

func Test_sendResultToPlugin(t *testing.T) {
	pChan := make(chan m2.KubeAuditResults)
	cChan := make(chan bool)
	auditTests := make([]*models.SubCategory, 0)
	ab := make([]*models.AuditBench, 0)
	ats := &models.AuditBench{Name: "bbb", TestSucceed: true}
	atf := &models.AuditBench{Name: "ccc", TestSucceed: false}
	ab = append(ab, ats)
	ab = append(ab, atf)
	mst := &models.SubCategory{Name: "aaa", AuditTests: ab}
	auditTests = append(auditTests, mst)
	go func() {
		<-pChan
		cChan <- true
	}()
	sendResultToPlugin(pChan, cChan, auditTests)

}
func Test_calculateFinalTotal(t *testing.T) {
	att := make([]models.AuditTestTotals, 0)
	atOne := models.AuditTestTotals{Fail: 2, Pass: 3, Warn: 1}
	atTwo := models.AuditTestTotals{Fail: 1, Pass: 5, Warn: 7}
	att = append(att, atOne)
	att = append(att, atTwo)
	res := calculateFinalTotal(att)
	assert.Equal(t, res.Warn, 8)
	assert.Equal(t, res.Pass, 8)
	assert.Equal(t, res.Fail, 3)
	str := printFinalResults([]models.AuditTestTotals{res})
	assert.Equal(t, str, "\u001B[34mTest Result Total\u001B[0m \u001B[32mPass:\u001B[0m 8 , \u001B[33mWarn:\u001B[0m 8 , \u001B[31mFail:\u001B[0m 3 ")
}
