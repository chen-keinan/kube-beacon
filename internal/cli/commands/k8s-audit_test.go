package commands

import (
	"fmt"
	m3 "github.com/chen-keinan/beacon/internal/cli/mocks"
	"github.com/chen-keinan/beacon/internal/logger"
 	"github.com/chen-keinan/beacon/internal/models"
 	m2 "github.com/chen-keinan/beacon/pkg/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/chen-keinan/go-command-eval/eval"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)



func TestRunAuditTests(t *testing.T) {
	tests := []struct {
		name              string
		testFile          string
		completedChan     chan bool
		plChan            chan m2.KubeAuditResults
		wantTestSucceeded bool
	}{

		{name: "Test_MultiCommandParams_OK", testFile: "CheckMultiParamOK.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: true},
		{name: "Test_MultiCommandParams_OK_With_IN", testFile: "CheckMultiParamOKWithIN.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: true},
		{name: "Test_MultiCommandParams_NOKWith_IN", testFile: "CheckMultiParamNOKWithIN.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: false},
		{name: "Test_MultiCommandParamsPass1stResultToNext", testFile: "CheckMultiParamPass1stResultToNext.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: false},
		{name: "Test_MultiCommandParamsComplex", testFile: "CheckMultiParamComplex.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: true},
		{name: "Test_MultiCommandParamsComplexOppositeEmptyReturn", testFile: "CheckInClauseOppositeEmptyReturn.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: false},
		{name: "Test_MultiCommandParamsComplexOppositeWithNumber", testFile: "CheckInClauseOppositeWithNum.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: false},
		{name: "Test_MultiCommand4_2_13", testFile: "CheckInClause4.2.13.yml", completedChan: make(chan bool), plChan: make(chan m2.KubeAuditResults), wantTestSucceeded: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := models.Audit{}
			err := yaml.Unmarshal(readTestData(tt.testFile, t), &ab)
			if err != nil {
				t.Errorf("failed to Unmarshal test file %s error : %s", tt.testFile, err.Error())
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			evalCmd := m3.NewMockCmdEvaluator(ctrl)
			testBench := ab.Categories[0].SubCategory.AuditTests[0]
			evalCmd.EXPECT().EvalCommand(testBench.AuditCommand, testBench.EvalExpr).Return(eval.CmdEvalResult{Match: tt.wantTestSucceeded, Error: nil}).Times(1)
			kb := K8sAudit{Evaluator: evalCmd, ResultProcessor: GetResultProcessingFunction([]string{}), PlChan: tt.plChan, CompletedChan: tt.completedChan}
			kb.runAuditTest(ab.Categories[0].SubCategory.AuditTests[0])
			assert.Equal(t, ab.Categories[0].SubCategory.AuditTests[0].TestSucceed, tt.wantTestSucceeded)
			go func() {
				<-tt.plChan
				tt.completedChan <- true
			}()
		})
	}
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
	evaluator := eval.NewEvalCmd()
	ka := NewK8sAudit(args, plChan, completedChan, []utils.FilesInfo{}, logger.GetLog(),evaluator)
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
	evaluator := eval.NewEvalCmd()
	ka := NewK8sAudit(args, plChan, completedChan, []utils.FilesInfo{}, logger.GetLog(),evaluator)
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
	fm := reportResultProcessor(ad, true)
	assert.True(t, len(fm) == 0)
	fm = reportResultProcessor(ad, false)
	assert.True(t, len(fm) == 1)
	assert.Equal(t, fm[0].Name, "1.2.1 aaa")
}

//Test_K8sSynopsis test
func Test_K8sSynopsis(t *testing.T) {
	args := []string{"a", "i=1.2.3"}
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	evaluator := eval.NewEvalCmd()
	ka := NewK8sAudit(args, plChan, completedChan, []utils.FilesInfo{}, logger.GetLog(),evaluator)
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
