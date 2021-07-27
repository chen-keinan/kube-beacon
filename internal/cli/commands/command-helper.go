package commands

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/pkg/filters"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/chen-keinan/beacon/ui"
	"github.com/mitchellh/colorstring"
	"gopkg.in/yaml.v2"
	"strings"
)

func printTestResults(at []*models.AuditBench, log *logger.BLogger) models.AuditTestTotals {
	var (
		warnCounter int
		passCounter int
		failCounter int
	)
	for _, a := range at {
		if a.NonApplicable {
			warnTest := colorstring.Color("[yellow][Warn]")
			log.Console(fmt.Sprintf("%s %s\n", warnTest, a.Name))
			warnCounter++
			continue
		}
		if a.TestSucceed {
			passTest := colorstring.Color("[green][Pass]")
			log.Console(fmt.Sprintf("%s %s\n", passTest, a.Name))
			passCounter++
		} else {
			failTest := colorstring.Color("[red][Fail]")
			log.Console(fmt.Sprintf("%s %s\n", failTest, a.Name))
			failCounter++
		}
	}
	return models.AuditTestTotals{Fail: failCounter, Pass: passCounter, Warn: warnCounter}
}

//AddFailedMessages add failed audit test to report data
func AddFailedMessages(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	av := make([]*models.AuditBench, 0)
	testSucceeded := NumFailedTest == 0
	at.TestSucceed = testSucceeded
	if !testSucceeded || at.NonApplicable {
		av = append(av, at)
	}
	return av
}

//AddAllMessages add all audit test to report data
func AddAllMessages(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	av := make([]*models.AuditBench, 0)
	testSucceeded := NumFailedTest == 0
	at.TestSucceed = testSucceeded
	av = append(av, at)
	return av
}

//TestLoader load tests from filesystem
//command-helper.go
//go:generate mockgen -destination=../../mocks/mock_TestLoader.go -package=mocks . TestLoader
type TestLoader interface {
	LoadAuditTests(fi []utils.FilesInfo) []*models.SubCategory
}

//AuditTestLoader object
type AuditTestLoader struct {
}

//NewFileLoader create new file loader
func NewFileLoader() TestLoader {
	return &AuditTestLoader{}
}

//LoadAuditTests load audit test from benchmark folder
func (tl AuditTestLoader) LoadAuditTests(auditFiles []utils.FilesInfo) []*models.SubCategory {
	auditTests := make([]*models.SubCategory, 0)
	audit := models.Audit{}
	for _, auditFile := range auditFiles {
		err := yaml.Unmarshal([]byte(auditFile.Data), &audit)
		if err != nil {
			panic("Failed to unmarshal audit test yaml file")
		}
		auditTests = append(auditTests, audit.Categories[0].SubCategory)
	}
	return auditTests
}

//FilterAuditTests filter audit tests by predicate chain
func FilterAuditTests(predicates []filters.Predicate, predicateParams []string, at *models.SubCategory) *models.SubCategory {
	return RunPredicateChain(predicates, predicateParams, len(predicates), at)
}

//RunPredicateChain call every predicate in chain and filter tests
func RunPredicateChain(predicates []filters.Predicate, predicateParams []string, size int, at *models.SubCategory) *models.SubCategory {
	if size == 0 {
		return at
	}
	return RunPredicateChain(predicates[1:size], predicateParams[1:size], size-1, predicates[size-1](at, predicateParams[size-1]))
}

// check weather are exist in array of specificTests
func isArgsExist(args []string, name string) bool {
	for _, n := range args {
		if n == name {
			return true
		}
	}
	return false
}

//GetResultProcessingFunction return processing function by specificTests
func GetResultProcessingFunction(args []string) ResultProcessor {
	if isArgsExist(args, common.Report) || isArgsExist(args, common.ReportFull) {
		return reportResultProcessor
	}
	return simpleResultProcessor
}

//getOutPutGeneratorFunction return output generator function
func getOutputGeneratorFunction(args []string) ui.OutputGenerator {
	if isArgsExist(args, common.Report) || isArgsExist(args, common.ReportFull) {
		return ReportOutputGenerator
	}
	return ConsoleOutputGenerator
}

//buildPredicateChain build chain of filters based on command criteria
func buildPredicateChain(args []string) []filters.Predicate {
	pc := make([]filters.Predicate, 0)
	for _, n := range args {
		switch {
		case strings.HasPrefix(n, common.IncludeParam):
			pc = append(pc, filters.IncludeAudit)
		case strings.HasPrefix(n, common.ExcludeParam):
			pc = append(pc, filters.ExcludeAudit)
		case strings.HasPrefix(n, common.NodeParam):
			pc = append(pc, filters.NodeAudit)
		case n == "a":
			pc = append(pc, filters.Basic)
		}
	}
	return pc
}

//buildPredicateParams build chain of filters params based on command criteria
func buildPredicateChainParams(args []string) []string {
	pp := make([]string, 0)
	pp = append(pp, args...)
	return pp
}

func filteredAuditBenchTests(auditTests []*models.SubCategory, pc []filters.Predicate, pp []string) []*models.SubCategory {
	ft := make([]*models.SubCategory, 0)
	for _, adt := range auditTests {
		filteredAudit := FilterAuditTests(pc, pp, adt)
		if len(filteredAudit.AuditTests) == 0 {
			continue
		}
		ft = append(ft, filteredAudit)
	}
	return ft
}

func executeTests(ft []*models.SubCategory, execTestFunc func(ad *models.AuditBench) []*models.AuditBench, log *logger.BLogger) []*models.SubCategory {
	completedTest := make([]*models.SubCategory, 0)
	log.Console(ui.K8sAuditTest)
	for _, f := range ft {
		tr := ui.ShowProgressBar(f, execTestFunc, log)
		completedTest = append(completedTest, tr)
	}
	return completedTest
}
