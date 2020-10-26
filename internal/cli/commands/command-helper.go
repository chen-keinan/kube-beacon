package commands

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/pkg/filters"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/chen-keinan/beacon/ui"
	"github.com/mitchellh/colorstring"
	"gopkg.in/yaml.v2"
	"strings"
)

func printTestResults(at []*models.AuditBench) {
	for _, a := range at {
		if a.TestSucceed {
			pass := colorstring.Color("[green][Pass]")
			log.Console(fmt.Sprintf("%s %s\n", pass, a.Name))
		} else {
			fail := colorstring.Color("[red][Fail]")
			log.Console(fmt.Sprintf("%s %s\n", fail, a.Name))
		}
	}
}

//AddFailedMessages add failed audit test to report data
func AddFailedMessages(at *models.AuditBench, NumFailedTest int) []*models.AuditBench {
	av := make([]*models.AuditBench, 0)
	testSucceeded := NumFailedTest == 0
	at.TestSucceed = testSucceeded
	if !testSucceeded {
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

//LoadAuditTests load audit test from benchmark folder
func LoadAuditTests() [][]*models.AuditBench {
	auditTests := make([][]*models.AuditBench, 0)
	audit := models.Audit{}
	auditFiles, err := utils.GetK8sBenchAuditFiles()
	if err != nil {
		panic(fmt.Sprintf("failed to read audit files %s", err))
	}
	for _, auditFile := range auditFiles {
		err := yaml.Unmarshal([]byte(auditFile.Data), &audit)
		if err != nil {
			panic("Failed to unmarshal audit test yaml file")
		}
		audit.Categories[0].SubCategory.AuditTests[0].Category = audit.Categories[0].SubCategory.Name
		auditTests = append(auditTests, audit.Categories[0].SubCategory.AuditTests)
	}
	return auditTests
}

//FilterAuditTests filter audit tests by predicate chain
func FilterAuditTests(predicates []filters.Predicate, predicateParams []string, at []*models.AuditBench) []*models.AuditBench {
	return RunPredicateChain(predicates, predicateParams, len(predicates), at)
}

//RunPredicateChain call every predicate in chain and filter tests
func RunPredicateChain(predicates []filters.Predicate, predicateParams []string, size int, at []*models.AuditBench) []*models.AuditBench {
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

//getResultProcessingFunction return processing function by specificTests
func getResultProcessingFunction(args []string) ResultProcessor {
	if isArgsExist(args, common.Report) {
		return reportResultProcessor
	}
	return simpleResultProcessor
}

//getOutPutGeneratorFunction return output generator function
func getOutputGeneratorFunction(args []string) ui.OutputGenerator {
	if isArgsExist(args, common.Report) {
		return reportOutputGenerator
	}
	return consoleOutputGenerator
}

//buildPredicateChain build chain of filters based on command criteria
func buildPredicateChain(args []string) []filters.Predicate {
	pc := make([]filters.Predicate, 0)
	for _, n := range args {
		switch {
		case strings.HasPrefix(n, common.IncludeParam):
			pc = append(pc, filters.IncludeAuditTest)
		case strings.HasPrefix(n, common.ExcludeParam):
			pc = append(pc, filters.ExcludeAuditTest)
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
