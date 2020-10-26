package filters

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"strings"
)

// Predicate filter audit tests cmd criteria
type Predicate func(tests []*models.AuditBench, params string) []*models.AuditBench

// IncludeAuditTest include audit tests , only included tests will be executed
var IncludeAuditTest Predicate = func(tests []*models.AuditBench, params string) []*models.AuditBench {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("i", params)
	for _, at := range tests {
		for _, sp := range spt {
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

// ExcludeAuditTest audit test from been executed
var ExcludeAuditTest Predicate = func(tests []*models.AuditBench, params string) []*models.AuditBench {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("e", params)
	for _, at := range tests {
		var skipTest bool
		for _, sp := range spt {
			if strings.Contains(at.Name, sp) {
				skipTest = true
			}
		}
		if skipTest {
			continue
		}
		sat = append(sat, at)
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
