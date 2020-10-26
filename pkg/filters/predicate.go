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

// ExcludeAuditTest audit test from been executed
var ExcludeAuditTest Predicate = func(tests []*models.AuditBench, params string) []*models.AuditBench {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("e", params)
	for _, sp := range spt {
		for _, at := range tests {
			if strings.Contains(at.Name, sp) {
				continue
			}
			sat = append(sat, at)
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
