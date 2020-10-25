package ui

import (
	"github.com/cheggaaa/pb"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"time"
)

var log = logger.GetLog()

// OutputGenerator for  audit results
type OutputGenerator func(at [][]*models.AuditBench)

//PrintOutput print audit test result to console
func PrintOutput(auditTests [][]*models.AuditBench, outputGenerator OutputGenerator) {
	log.Console(auditResult)
	outputGenerator(auditTests)
}

//ShowProgressBar execute audit test and show progress bar
func ShowProgressBar(a []*models.AuditBench, f func(ad *models.AuditBench) []*models.AuditBench) []*models.AuditBench {
	completedTest := make([]*models.AuditBench, 0)
	log.Console(a[0].Category)
	bar := pb.StartNew(len(a))
	for _, test := range a {
		ar := f(test)
		completedTest = append(completedTest, ar...)
		bar.Increment()
		time.Sleep(time.Millisecond * 20)
	}
	bar.Finish()
	return completedTest
}
