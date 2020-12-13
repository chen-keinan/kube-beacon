package commands

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/reports"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/chen-keinan/beacon/pkg/filters"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/chen-keinan/beacon/ui"
	"github.com/mitchellh/colorstring"
	"strconv"
	"strings"
)

var log = logger.GetLog()

//K8sAudit k8s benchmark object
type K8sAudit struct {
	Command         shell.Executor
	ResultProcessor ResultProcessor
	OutputGenerator ui.OutputGenerator
	FileLoader      TestLoader
	PredicateChain  []filters.Predicate
	PredicateParams []string
	Spec            string
	Version         string
}

// ResultProcessor process audit results
type ResultProcessor func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench

// ConsoleOutputGenerator print audit tests to stdout
var ConsoleOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory) {
	granTotal := make([]models.AuditTestTotals, 0)
	for _, a := range at {
		log.Console(fmt.Sprintf("%s %s\n", "[Category]", a.Name))
		categoryTotal := printTestResults(a.AuditTests)
		granTotal = append(granTotal, categoryTotal)
	}
	finalTotal := calculateFinalTotal(granTotal)
	passTest := colorstring.Color("[green]Pass:")
	failTest := colorstring.Color("[red]Fail:")
	warnTest := colorstring.Color("[yellow]Warn:")
	title := colorstring.Color("[blue]Test Result Total")
	log.Console(fmt.Sprintf("%s %s %d , %s %d , %s %d ", title, passTest, finalTotal.Pass, warnTest, finalTotal.Warn, failTest, finalTotal.Fail))
}

func calculateFinalTotal(granTotal []models.AuditTestTotals) models.AuditTestTotals {
	var (
		warn int
		fail int
		pass int
	)
	for _, total := range granTotal {
		warn = warn + total.Warn
		fail = fail + total.Fail
		pass = pass + total.Pass
	}
	return models.AuditTestTotals{Pass: pass, Fail: fail, Warn: warn}
}

// ReportOutputGenerator print failed audit test to human report
var ReportOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory) {
	for _, a := range at {
		log.Table(reports.GenerateAuditReport(a.AuditTests))
	}
}

// simpleResultProcessor process audit results to stdout print only
var simpleResultProcessor ResultProcessor = func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	return AddAllMessages(at, NumFailedTest)
}

// ResultProcessor process audit results to std out and failure results
var reportResultProcessor ResultProcessor = func(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	// append failed messages
	return AddFailedMessages(at, NumFailedTest)
}

//NewK8sAudit new audit object
func NewK8sAudit(args []string, spec, version string) *K8sAudit {
	return &K8sAudit{Command: shell.NewShellExec(),
		PredicateChain:  buildPredicateChain(args),
		PredicateParams: buildPredicateChainParams(args),
		ResultProcessor: GetResultProcessingFunction(args),
		OutputGenerator: getOutputGeneratorFunction(args),
		FileLoader:      NewFileLoader(),
		Spec:            spec,
		Version:         version}

}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return startup.GetHelpSynopsis()
}

//Run execute the full k8s benchmark
func (bk *K8sAudit) Run(args []string) int {
	// load audit tests fro benchmark folder
	auditTests := bk.FileLoader.LoadAuditTests(bk.Spec, bk.Version)
	// filter tests by cmd criteria
	ft := filteredAuditBenchTests(auditTests, bk.PredicateChain, bk.PredicateParams)
	//execute audit tests and show it in progress bar
	completedTest := executeTests(ft, bk.runAuditTest)
	// generate output data
	ui.PrintOutput(completedTest, bk.OutputGenerator)
	return 0
}

// runAuditTest execute category of audit tests
func (bk *K8sAudit) runAuditTest(at *models.AuditBench) []*models.AuditBench {
	auditRes := make([]*models.AuditBench, 0)
	if at.NonApplicable {
		auditRes = append(auditRes, at)
		return auditRes
	}
	cmdTotalRes := make([]string, 0)
	// execute audit test command
	for index := range at.AuditCommand {
		res := bk.execCommand(at, index, cmdTotalRes, make([]IndexValue, 0))
		cmdTotalRes = append(cmdTotalRes, res)
	}
	// evaluate command result with expression
	NumFailedTest := bk.evalExpression(at, cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0)
	// continue with result processing
	auditRes = append(auditRes, bk.ResultProcessor(at, NumFailedTest)...)
	return auditRes
}

func (bk *K8sAudit) addDummyCommandResponse(expr string, index int, n string) string {
	if n == "[^\"]\\S*'\n" || n == "" || n == common.EmptyValue {
		spExpr := utils.SeparateExpr(expr)
		for _, expr := range spExpr {
			if expr.Type == common.SingleValue {
				if !strings.Contains(expr.Expr, fmt.Sprintf("'$%d'", index)) {
					if strings.Contains(expr.Expr, fmt.Sprintf("$%d", index)) {
						return common.NotValidNumber
					}
				}
			}
		}
		return common.EmptyValue
	}
	return n
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
			paramNum, err := strconv.Atoi(param)
			if err != nil {
				log.Console(fmt.Sprintf("failed to convert param for command %s", cmd))
				continue
			}
			if paramNum < len(prevResult) {
				n := bk.addDummyCommandResponse(at.EvalExpr, index, prevResult[paramNum])
				newRes = append(newRes, IndexValue{index: paramNum, value: n})
			}
		}
		commandRes := bk.execCmdWithParams(newRes, len(newRes), make([]IndexValue, 0), cmd, make([]string, 0))
		sb := strings.Builder{}
		for _, cr := range commandRes {
			sb.WriteString(utils.AddNewLineToNonEmptyStr(cr))
		}
		return sb.String()
	}
	result, _ := bk.Command.Exec(cmd)
	if result.Stderr != "" {
		log.Console(fmt.Sprintf("Failed to execute command %s\n %s", result.Stderr, cmd))
	}
	return bk.addDummyCommandResponse(at.EvalExpr, index, result.Stdout)
}

func (bk *K8sAudit) execCmdWithParams(arr []IndexValue, index int, prevResHolder []IndexValue, currCommand string, resArr []string) []string {
	if len(arr) == 0 {
		return execShellCmd(prevResHolder, resArr, currCommand, bk.Command)
	}
	sArr := strings.Split(utils.RemoveNewLineSuffix(arr[0].value), "\n")
	for _, a := range sArr {
		prevResHolder = append(prevResHolder, IndexValue{index: arr[0].index, value: a})
		resArr = bk.execCmdWithParams(arr[1:index], index-1, prevResHolder, currCommand, resArr)
		prevResHolder = prevResHolder[:len(prevResHolder)-1]
	}
	return resArr
}

func execShellCmd(prevResHolder []IndexValue, resArr []string, currCommand string, se shell.Executor) []string {
	for _, param := range prevResHolder {
		if param.value == common.EmptyValue || param.value == common.NotValidNumber || param.value == "" {
			resArr = append(resArr, param.value)
			break
		}
		cmd := strings.ReplaceAll(currCommand, fmt.Sprintf("#%d", param.index), param.value)
		result, _ := se.Exec(cmd)
		if result.Stderr != "" {
			log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
		}
		if len(strings.TrimSpace(result.Stdout)) == 0 {
			result.Stdout = common.EmptyValue
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
	outputs := strings.Split(utils.RemoveNewLineSuffix(commandRes[0]), "\n")
	for _, o := range outputs {
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
	testSucceeded, err := evalCommandExpr(strings.ReplaceAll(expr, common.EmptyValue, ""))
	if err != nil {
		log.Console(fmt.Sprintf("failed to evaluate command expr %s for audit test %s", expr, at.Name))
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
