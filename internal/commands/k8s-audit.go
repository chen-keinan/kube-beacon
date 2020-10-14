package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s/models"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
	"strings"
)

//ValidateExprData expr data
type ValidateExprData struct {
	index     int
	resultArr []string
	atb       *models.AuditBench
	origSize  int
	Total     int
	Match     int
}

//NextValidExprData return the next recursive ValidExprData
func (ve ValidateExprData) NextValidExprData() ValidateExprData {
	return ValidateExprData{resultArr: ve.resultArr[1:ve.index], index: ve.index - 1, atb: ve.atb, origSize: ve.origSize}
}

// NewValidExprData return new instance of ValidExprData
func NewValidExprData(arr []string, at *models.AuditBench) ValidateExprData {
	return ValidateExprData{resultArr: arr, index: len(arr), atb: at, origSize: len(arr)}
}

//K8sAudit k8s benchmark object
type K8sAudit struct {
}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return "-a , --audit run benchmark audit tests"
}

//Run execute benchmark command
func (bk K8sAudit) Run(args []string) int {
	audit := models.Audit{}
	auditFiles, err := utils.GetK8sBenchAuditFiles()
	if err != nil {
		panic(fmt.Sprintf("failed to read audit files %s", err))
	}
	for _, auditFile := range auditFiles {
		err := json.Unmarshal([]byte(auditFile.Data), &audit)
		if err != nil {
			panic("Failed to unmarshal audit test json file")
		}
		for _, ac := range audit.Categories {
			bk.runTests(ac)
		}
	}
	return 0
}

func (bk K8sAudit) runTests(ac models.Category) {
	for _, at := range ac.SubCategory.AuditTests {
		resArr := make([]string, 0)
		for _, val := range at.AuditCommand {
			result, err := shell.NewShellExec().Exec(val)
			if err != nil {
				fmt.Printf("Failed to execute command %s", err.Error())
				continue
			}
			resArr = append(resArr, result.Stdout)
		}
		data := NewValidExprData(resArr, at)
		bk.evalExpression(data, make([]string, 0))
		bk.printTestResults(data.atb)
	}
}

func (bk K8sAudit) printTestResults(at *models.AuditBench) {
	if at.TestResult.NumOfSuccess == at.TestResult.NumOfExec {
		fmt.Print(emoji.Sprintf(":check_mark_button: %s\n", at.Name))
	} else {
		fmt.Print(emoji.Sprintf(":cross_mark: %s\n", at.Name))

	}
}

func (bk K8sAudit) evalExpression(ved ValidateExprData, combArr []string) {
	if len(ved.resultArr) == 0 {
		return
	}
	outputs := strings.Split(ved.resultArr[0], "\n")
	for _, o := range outputs {
		if len(o) == 0 && len(outputs) > 1 {
			continue
		}
		combArr = append(combArr, o)
		bk.evalExpression(ved.NextValidExprData(), combArr)
		if ved.origSize == len(combArr) {
			expr := ved.atb.Sanitize(combArr, ved.atb.EvalExpr)
			ved.atb.TestResult.NumOfExec++
			count, err := bk.evalCommandExpr(ved.atb, expr)
			if err != nil {
				fmt.Println(err)
			}
			ved.atb.TestResult.NumOfSuccess += count
		}
		combArr = combArr[:len(combArr)-1]
	}

}

func (bk K8sAudit) evalCommandExpr(at *models.AuditBench, expr string) (int, error) {

	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to build evaluation command expr for\n %s", at.Name)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate command expr for audit test %s", at.Name)
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}

//Synopsis for help
func (bk K8sAudit) Synopsis() string {
	return bk.Help()
}
