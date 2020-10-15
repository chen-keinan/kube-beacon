package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s/models"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
	"strconv"
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
	Command shell.Executor
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
		for index, val := range at.AuditCommand {
			cmd := bk.UpdateCommand(at, index, val, resArr)
			if cmd == "" {
				continue
			}
			result, _ := bk.Command.Exec(cmd)
			if result.Stderr != "" {
				resArr = append(resArr, "")
				fmt.Printf("Failed to execute command %s", result.Stderr)
				continue
			}
			resArr = append(resArr, result.Stdout)
		}
		data := NewValidExprData(resArr, at)
		if len(at.AuditCommand) == len(resArr) {
			bk.evalExpression(data, make([]string, 0))
		} else {
			at.TestResult.NumOfExec = 1
			at.TestResult.NumOfSuccess = 0
		}
		bk.printTestResults(data.atb)
	}
}

//UpdateCommand update the cmd command with params values
func (bk K8sAudit) UpdateCommand(at *models.AuditBench, index int, val string, resArr []string) string {
	params := at.CommandParams[index]
	if len(params) > 0 {
		for _, param := range params {
			x, err := strconv.Atoi(param)
			if err != nil {
				fmt.Printf("failed to translate param %s to number", param)
				continue
			}
			n := resArr[x]
			switch {
			case n == "[^\"]\\S*'\n" || n == "":
				return ""
			case strings.Contains(n, "\n"):
				nl := n[len(n)-1:]
				if nl == "\n" {
					n = strings.Trim(n, "\n")
				}
			}
			return strings.ReplaceAll(val, fmt.Sprintf("#%d", x), n)
		}
	}
	return val
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
