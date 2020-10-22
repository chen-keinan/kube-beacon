package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/reports"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/chen-keinan/beacon/pkg/utils"
	"strconv"
	"strings"
)

var log = logger.GetLog()

//ValidateExprData expr data
type ValidateExprData struct {
	index     int
	resultArr []string
	atb       *models.AuditBench
}

//NextValidExprData return the next recursive ValidExprData
func (ve ValidateExprData) NextValidExprData() ValidateExprData {
	return ValidateExprData{resultArr: ve.resultArr[1:ve.index], index: ve.index - 1, atb: ve.atb}
}

// NewValidExprData return new instance of ValidExprData
func NewValidExprData(arr []string, at *models.AuditBench) ValidateExprData {
	return ValidateExprData{resultArr: arr, index: len(arr), atb: at}
}

//K8sAudit k8s benchmark object
type K8sAudit struct {
	Command         shell.Executor
	specificTests   []string
	resultProcessor ResultProcessor
}

// ResultProcessor process audit results
type ResultProcessor func(vd ValidateExprData, NumFailedTest int) []*models.AuditBench

// simpleResultProcessor process audit results to stdout print only
var simpleResultProcessor ResultProcessor = func(vd ValidateExprData, NumFailedTest int) []*models.AuditBench {
	printTestResults(vd.atb, NumFailedTest)
	return []*models.AuditBench{}
}

// ResultProcessor process audit results to std out and failure results
var reportResultProcessor ResultProcessor = func(vd ValidateExprData, NumFailedTest int) []*models.AuditBench {
	// append failed messages
	return AddFailedMessages(vd, NumFailedTest)
}

//NewK8sAudit new audit object
func NewK8sAudit(args []string) *K8sAudit {
	return &K8sAudit{Command: shell.NewShellExec(), specificTests: getSpecificTestsToExecute(args), resultProcessor: getResultProcessingFunction(args)}
}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return startup.GetHelpSynopsis()
}

//Run execute benchmark command
func (bk *K8sAudit) Run(args []string) int {
	auditRes := make([]*models.AuditBench, 0)
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
			ar := bk.runTests(ac)
			auditRes = append(auditRes, ar...)
		}
	}
	// generate report data
	reports.GenerateAuditReport(auditRes)
	return 0
}

func (bk *K8sAudit) runTests(ac models.Category) []*models.AuditBench {
	auditRes := make([]*models.AuditBench, 0)
	for _, at := range ac.SubCategory.AuditTests {
		if utils.ExcludeAuditTest(bk.specificTests, at.Name) {
			continue
		}
		cmdTotalRes := make([]string, 0)
		for index := range at.AuditCommand {
			res := bk.execCommand(at, index, cmdTotalRes, make([]IndexValue, 0))
			cmdTotalRes = append(cmdTotalRes, res)
		}
		data := NewValidExprData(cmdTotalRes, at)
		// evaluate command result with expression
		NumFailedTest := bk.evalExpression(at, cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0)
		// continue with result processing
		auditRes = append(auditRes, bk.resultProcessor(data, NumFailedTest)...)
	}
	return auditRes
}

func (bk *K8sAudit) addDummyCommandResponse(at *models.AuditBench, index int) string {
	spExpr := utils.SeparateExpr(at.EvalExpr)
	for _, expr := range spExpr {
		if expr.Type == common.SingleValue {
			if !strings.Contains(expr.Expr, fmt.Sprintf("'$%d'", index)) {
				if strings.Contains(expr.Expr, fmt.Sprintf("$%d", index)) {
					return common.NotValidNumber
				}
			}
		}
	}
	return common.NotValidString
}

//IndexValue hold command index and result
type IndexValue struct {
	index int
	value string
}

func (bk *K8sAudit) execCommand(at *models.AuditBench, index int, prevResult []string, newRes []IndexValue) string {
	cmd := at.AuditCommand[index]
	paramArr, ok := at.CommandParams[index]
	if ok {
		for _, param := range paramArr {
			parmNum, err := strconv.Atoi(param)
			if err != nil {
				log.Console(fmt.Sprintf("failed to convert param for command %s", cmd))
				continue
			}
			if parmNum < len(prevResult) {
				n := prevResult[parmNum]
				if n == "[^\"]\\S*'\n" || n == "" || n == common.NotValidString {
					n = bk.addDummyCommandResponse(at, index)
				}
				newRes = append(newRes, IndexValue{index: parmNum, value: n})
			}
		}
		commandRes := bk.execCmdWithParams(newRes, len(newRes), make([]IndexValue, 0), cmd, make([]string, 0))
		sb := strings.Builder{}
		for _, cr := range commandRes {
			sb.WriteString(fmt.Sprintf("%s\n", cr))
		}
		return sb.String()
	}
	result, _ := bk.Command.Exec(cmd)
	if result.Stderr != "" {
		log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
	}
	return result.Stdout

}

func (bk *K8sAudit) execCmdWithParams(arr []IndexValue, index int, prevResHolder []IndexValue, currCommand string, resArr []string) []string {
	if len(arr) == 0 {
		return execShellCmd(prevResHolder, resArr, currCommand, bk.Command)
	}
	sArr := strings.Split(arr[0].value, "\n")
	for _, a := range sArr {
		prevResHolder = append(prevResHolder, IndexValue{index: arr[0].index, value: a})
		resArr = bk.execCmdWithParams(arr[1:index], index-1, prevResHolder, currCommand, resArr)
		prevResHolder = prevResHolder[:len(prevResHolder)-1]
	}
	return resArr
}

func execShellCmd(prevResHolder []IndexValue, resArr []string, currCommand string, se shell.Executor) []string {
	for _, param := range prevResHolder {
		if param.value == common.NotValidString || param.value == common.NotValidNumber || param.value == "" {
			resArr = append(resArr, param.value)
			break
		}
		cmd := strings.ReplaceAll(currCommand, fmt.Sprintf("#%d", param.index), param.value)
		result, _ := se.Exec(cmd)
		if result.Stderr != "" {
			resArr = append(resArr, "")
			log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
		}
		resArr = append(resArr, result.Stdout)
	}
	return resArr
}

//evalExpression expression eval as cartesian product
func (bk *K8sAudit) evalExpression(at *models.AuditBench,
	commandRes []string, commResSize int, permutationArr []string, testFailure int) int {
	if len(commandRes) == 0 {
		return evalCommand(at, permutationArr, testFailure)
	}
	outputs := strings.Split(commandRes[0], "\n")
	for _, o := range outputs {
		if len(o) == 0 && len(outputs) > 1 {
			continue
		}
		permutationArr = append(permutationArr, o)
		testFailure = bk.evalExpression(at, commandRes[1:commResSize], commResSize-1, permutationArr, testFailure)
		permutationArr = permutationArr[:len(permutationArr)-1]
	}
	return testFailure
}

func evalCommand(at *models.AuditBench, permutationArr []string, testExec int) int {
	// build command expression with params
	expr := at.CmdExprBuilder(permutationArr, at.EvalExpr)
	testExec++
	// eval command expression
	testSucceeded, err := evalCommandExpr(expr)
	if err != nil {
		log.Console(fmt.Sprintf("failed to evaluate command expr for audit test %s", at.Name))
	}
	return testExec - testSucceeded
}

func evalCommandExpr(expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}

//Synopsis for help
func (bk *K8sAudit) Synopsis() string {
	return bk.Help()
}
