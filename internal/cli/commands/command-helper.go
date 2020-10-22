package commands

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/kyokomi/emoji"
	"strings"
)

func printTestResults(at *models.AuditBench, NumFailedTest int) {
	testSucceeded := NumFailedTest == 0
	at.TestSucceed = testSucceeded
	if testSucceeded {
		log.Console(emoji.Sprintf(":check_mark_button: %s\n", at.Name))
	} else {
		log.Console(emoji.Sprintf(":cross_mark: %s\n", at.Name))
	}
	at.TestSucceed = testSucceeded
}

//AddFailedMessages add failed audit test to report data
func AddFailedMessages(data ValidateExprData, NumFailedTest int) []*models.AuditBench {
	av := make([]*models.AuditBench, 0)
	testSucceeded := NumFailedTest == 0
	data.atb.TestSucceed = testSucceeded
	if !testSucceeded {
		av = append(av, data.atb)
	}
	return av
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

//getSpecificTestsToExecute return processing function by specificTests
func getSpecificTestsToExecute(args []string) []string {
	for _, n := range args {
		if strings.HasPrefix(n, "s=") {
			values := strings.ReplaceAll(n, "s=", "")
			return strings.Split(values, ";")
		}
	}
	return []string{}
}
