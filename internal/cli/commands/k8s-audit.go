package commands

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/reports"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/chen-keinan/beacon/pkg/filters"
	m2 "github.com/chen-keinan/beacon/pkg/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/chen-keinan/beacon/ui"
	"github.com/chen-keinan/go-command-eval/eval"
	"github.com/mitchellh/colorstring"
)

//K8sAudit k8s benchmark object
type K8sAudit struct {
	ResultProcessor ResultProcessor
	OutputGenerator ui.OutputGenerator
	FileLoader      TestLoader
	PredicateChain  []filters.Predicate
	PredicateParams []string
	PlChan          chan m2.KubeAuditResults
	CompletedChan   chan bool
	FilesInfo       []utils.FilesInfo
	Evaluator       eval.CmdEvaluator
	Log             *logger.BLogger
}

// ResultProcessor process audit results
type ResultProcessor func(at *models.AuditBench, match bool) []*models.AuditBench

// ConsoleOutputGenerator print audit tests to stdout
var ConsoleOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.BLogger) {
	grandTotal := make([]models.AuditTestTotals, 0)
	for _, a := range at {
		log.Console(fmt.Sprintf("%s %s\n", "[Category]", a.Name))
		categoryTotal := printTestResults(a.AuditTests, log)
		grandTotal = append(grandTotal, categoryTotal)
	}
	log.Console(printFinalResults(grandTotal))
}

func printFinalResults(grandTotal []models.AuditTestTotals) string {
	finalTotal := calculateFinalTotal(grandTotal)
	passTest := colorstring.Color("[green]Pass:")
	failTest := colorstring.Color("[red]Fail:")
	warnTest := colorstring.Color("[yellow]Warn:")
	title := colorstring.Color("[blue]Test Result Total")
	return fmt.Sprintf("%s %s %d , %s %d , %s %d ", title, passTest, finalTotal.Pass, warnTest, finalTotal.Warn, failTest, finalTotal.Fail)
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
var ReportOutputGenerator ui.OutputGenerator = func(at []*models.SubCategory, log *logger.BLogger) {
	for _, a := range at {
		log.Table(reports.GenerateAuditReport(a.AuditTests))
	}
}

// simpleResultProcessor process audit results to stdout print only
var simpleResultProcessor ResultProcessor = func(at *models.AuditBench, match bool) []*models.AuditBench {
	return AddAllMessages(at, match)
}

// ResultProcessor process audit results to std out and failure results
var reportResultProcessor ResultProcessor = func(at *models.AuditBench, match bool) []*models.AuditBench {
	// append failed messages
	return AddFailedMessages(at, match)
}

//CmdEvaluator interface expose one method to evaluate command with evalExpr
//k8s-audit.go
//go:generate mockgen -destination=../mocks/mock_CmdEvaluator.go -package=mocks . CmdEvaluator
type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) eval.CmdEvalResult
}

//NewK8sAudit new audit object
func NewK8sAudit(filters []string, plChan chan m2.KubeAuditResults, completedChan chan bool, fi []utils.FilesInfo, log *logger.BLogger, evaluator CmdEvaluator) *K8sAudit {
	return &K8sAudit{
		PredicateChain:  buildPredicateChain(filters),
		PredicateParams: buildPredicateChainParams(filters),
		ResultProcessor: GetResultProcessingFunction(filters),
		OutputGenerator: getOutputGeneratorFunction(filters),
		FileLoader:      NewFileLoader(),
		PlChan:          plChan,
		FilesInfo:       fi,
		Evaluator:       evaluator,
		Log:             log,
		CompletedChan:   completedChan}
}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return startup.GetHelpSynopsis()
}

//Run execute the full k8s benchmark
func (bk *K8sAudit) Run(args []string) int {
	// load audit tests fro benchmark folder
	auditTests := bk.FileLoader.LoadAuditTests(bk.FilesInfo)
	// filter tests by cmd criteria
	ft := filteredAuditBenchTests(auditTests, bk.PredicateChain, bk.PredicateParams)
	//execute audit tests and show it in progress bar
	completedTest := executeTests(ft, bk.runAuditTest, bk.Log)
	// generate output data
	ui.PrintOutput(completedTest, bk.OutputGenerator, bk.Log)
	// send test results to plugin
	sendResultToPlugin(bk.PlChan, bk.CompletedChan, completedTest)
	return 0
}

func sendResultToPlugin(plChan chan m2.KubeAuditResults, completedChan chan bool, auditTests []*models.SubCategory) {
	ka := m2.KubeAuditResults{BenchmarkType: "k8s", Categories: make([]m2.AuditBenchResult, 0)}
	for _, at := range auditTests {
		for _, ab := range at.AuditTests {
			var testResult = "FAIL"
			if ab.TestSucceed {
				testResult = "PASS"
			}
			abr := m2.AuditBenchResult{Category: at.Name, ProfileApplicability: ab.ProfileApplicability, Description: ab.Description, AuditCommand: ab.AuditCommand, Remediation: ab.Remediation, Impact: ab.Impact, DefaultValue: ab.DefaultValue, References: ab.References, TestResult: testResult}
			ka.Categories = append(ka.Categories, abr)
		}
	}
	plChan <- ka
	<-completedChan
}

// runAuditTest execute category of audit tests
func (bk *K8sAudit) runAuditTest(at *models.AuditBench) []*models.AuditBench {
	auditRes := make([]*models.AuditBench, 0)
	if at.NonApplicable {
		auditRes = append(auditRes, at)
		return auditRes
	}
	// execute audit test command
	cmdEvalResult := bk.Evaluator.EvalCommand(at.AuditCommand, at.EvalExpr)
	// continue with result processing
	auditRes = append(auditRes, bk.ResultProcessor(at, cmdEvalResult.Match)...)
	return auditRes
}

//Synopsis for help
func (bk *K8sAudit) Synopsis() string {
	return bk.Help()
}
