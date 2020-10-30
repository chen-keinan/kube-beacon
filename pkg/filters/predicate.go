package filters

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"strings"
)

// Predicate filter audit tests cmd criteria
type Predicate func(tests *models.SubCategory, params string) *models.SubCategory

// IncludeAudit include audit tests , only included tests will be executed
var IncludeAudit Predicate = func(tests *models.SubCategory, params string) *models.SubCategory {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("i", params)
	// check if param include category
	for _, sp := range spt {
		if strings.HasPrefix(tests.Name, sp) {
			return tests
		}
	}
	// check tests
	for _, at := range tests.AuditTests {
		for _, sp := range spt {
			if strings.HasPrefix(at.Name, sp) {
				sat = append(sat, at)
			}
		}
	}
	if len(sat) == 0 {
		return &models.SubCategory{Name: tests.Name, AuditTests: make([]*models.AuditBench, 0)}
	}
	return &models.SubCategory{Name: tests.Name, AuditTests: sat}
}

// ExcludeAudit audit test from been executed
var ExcludeAudit Predicate = func(tests *models.SubCategory, params string) *models.SubCategory {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("e", params)
	// if exclude category
	for _, sp := range spt {
		if strings.HasPrefix(tests.Name, sp) {
			return &models.SubCategory{Name: tests.Name, AuditTests: []*models.AuditBench{}}
		}
	}
	for _, at := range tests.AuditTests {
		var skipTest bool
		for _, sp := range spt {
			if strings.HasPrefix(at.Name, sp) {
				skipTest = true
			}
		}
		if skipTest {
			continue
		}
		sat = append(sat, at)
	}
	return &models.SubCategory{Name: tests.Name, AuditTests: sat}
}

// NodeAudit audit test from been executed
var NodeAudit Predicate = func(tests *models.SubCategory, params string) *models.SubCategory {
	sat := make([]*models.AuditBench, 0)
	spt := utils.GetAuditTestsList("n", params)
	// check tests
	for _, at := range tests.AuditTests {
		for _, sp := range spt {
			if strings.ToLower(at.ProfileApplicability) == sp {
				sat = append(sat, at)
			}
		}
	}
	if len(sat) == 0 {
		return &models.SubCategory{Name: tests.Name, AuditTests: make([]*models.AuditBench, 0)}
	}
	return &models.SubCategory{Name: tests.Name, AuditTests: sat}
}

// Basic filter by specific audit tests as set in command
var Basic Predicate = func(tests *models.SubCategory, params string) *models.SubCategory {
	return tests
}
