package filters

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"strings"
)

// Predicate filter audit tests cmd criteria
type Predicate func(tests []*models.AuditBench, params string) []*models.AuditBench

// SpecificTest Basic test do not filter at all
var SpecificTest Predicate = func(tests []*models.AuditBench, params string) []*models.AuditBench {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetSpecificTestsToExecute(params)
	for _, sp := range spt {
		for _, at := range tests {
			if strings.Contains(at.Name, sp) {
				sat = append(sat, at)
			}
		}
	}
	if len(sat) == 0 {
		return tests
	}
	return sat
}

// Basic filter by specific audit tests as set in command
var Basic Predicate = func(tests []*models.AuditBench, params string) []*models.AuditBench {
	return tests
}
