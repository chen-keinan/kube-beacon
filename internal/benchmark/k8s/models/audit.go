package models

import (
	"encoding/json"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"strings"
)

//Audit data model
type Audit struct {
	BenchmarkType string     `json:"benchmark_type"`
	Categories    []Category `json:"categories"`
}

//Category data model
type Category struct {
	Name        string      `json:"name"`
	SubCategory SubCategory `json:"sub_category"`
}

//SubCategory data model
type SubCategory struct {
	Name       string        `json:"name"`
	AuditTests []*AuditBench `json:"audit_tests"`
}

//AuditBench data model
type AuditBench struct {
	Name                 string   `mapstructure:"name" json:"name"`
	ProfileApplicability string   `mapstructure:"profile_applicability" json:"profile_applicability"`
	Description          string   `mapstructure:"description" json:"description"`
	AuditCommand         []string `mapstructure:"audit" json:"audit"`
	CheckType            string   `mapstructure:"check_type" json:"check_type"`
	Remediation          string   `mapstructure:"remediation" json:"remediation"`
	Impact               string   `mapstructure:"impact" json:"impact"`
	DefaultValue         string   `mapstructure:"default_value" json:"default_value"`
	References           []string `mapstructure:"references" json:"references"`
	EvalExpr             string   `mapstructure:"eval_expr" json:"eval_expr"`
	Sanitize             utils.ExprSanitize
	TestResult           *AuditResult
	CommandParams        map[int][]string
}

//AuditResult data
type AuditResult struct {
	NumOfExec    int
	NumOfSuccess int
}

//UnmarshalJSON over unmarshall to add logic
func (at *AuditBench) UnmarshalJSON(data []byte) error {
	var res map[string]interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	err := mapstructure.Decode(res, &at)
	if err != nil {
		return err
	}
	switch at.CheckType {
	case "multi_param":
		at.Sanitize = utils.ExprSanitizeMultiProcessParam
	}
	at.TestResult = &AuditResult{}
	for index, command := range at.AuditCommand {
		at.CommandParams = make(map[int][]string)
		findIndex(command, "#", index, at.CommandParams)
	}
	return nil
}

// find all params in command to be replace with output
func findIndex(s, c string, commandIndex int, locations map[int][]string) {
	b := strings.Index(s, c)
	if b == -1 {
		return
	}
	if locations[commandIndex] == nil {
		locations[commandIndex] = make([]string, 0)
	}
	locations[commandIndex] = append(locations[commandIndex], s[b+1:b+2])
	findIndex(s[b+2:], c, commandIndex, locations)
}
