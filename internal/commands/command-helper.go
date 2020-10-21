package commands

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/kyokomi/emoji"
)

func printTestResults(at *models.AuditBench) {
	if at.TestResult.NumOfSuccess == at.TestResult.NumOfExec {
		log.Console(emoji.Sprintf(":check_mark_button: %s\n", at.Name))
	} else {
		log.Console(emoji.Sprintf(":cross_mark: %s\n", at.Name))
	}
}

//AddFailedMessages add failed audit test to report data
func AddFailedMessages(data ValidateExprData) []*models.AuditBench {
	av := make([]*models.AuditBench, 0)
	if data.atb.TestResult.NumOfSuccess != data.atb.TestResult.NumOfExec {
		av = append(av, data.atb)
	}
	return av
}

// check weather are exist in array of args
func isArgsExist(args []string, name string) bool {
	for _, n := range args {
		if n == name {
			return true
		}
	}
	return false
}

//getResultProcessingFunction return processing function by args
func getResultProcessingFunction(args []string) ResultProcessor {
	if isArgsExist(args, common.Report) {
		return reportResultProcessor
	}
	return simpleResultProcessor
}
