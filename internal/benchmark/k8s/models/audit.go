package models

import (
	"encoding/json"
	"github.com/chen-keinan/beacon/internal/common"
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
	Name       string       `json:"name"`
	AuditTests []AuditBench `json:"audit_tests"`
}

//AuditBench data model
type AuditBench struct {
	Name                 string   `mapstructure:"name" json:"name"`
	ProfileApplicability string   `mapstructure:"profile_applicability" json:"profile_applicability"`
	Description          string   `mapstructure:"description" json:"description"`
	AuditCommand         string   `mapstructure:"audit" json:"audit"`
	CheckType            string   `mapstructure:"check_type" json:"check_type"`
	Remediation          string   `mapstructure:"remediation" json:"remediation"`
	Impact               string   `mapstructure:"impact" json:"impact"`
	DefaultValue         string   `mapstructure:"default_value" json:"default_value"`
	References           []string `mapstructure:"references" json:"references"`
	EvalExpr             string   `mapstructure:"eval_expr" json:"eval_expr"`
	Sanitize             ExprSanitize
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
	case "ownership":
		at.Sanitize = exprSanitizeOwnership
	case "permission":
		at.Sanitize = exprSanitizePermission
	case "process_param":
		at.Sanitize = exprSanitizeProcessParam
	case "multi_process_param":
		at.Sanitize = exprSanitizeMultiProcessParam
	}
	return nil
}

//ExprSanitize sanitize expr
type ExprSanitize func(output, expr string) string

var exprSanitizeOwnership ExprSanitize = func(output, expr string) string {
	return SanitizeRegExOutPut(output, expr)
}

var exprSanitizeProcessParam ExprSanitize = func(output, expr string) string {
	return SanitizeRegExOutPut(output, expr)
}

var exprSanitizeMultiProcessParam ExprSanitize = func(output, expr string) string {
	var s string
	if strings.Contains(expr, common.GrepRegex) {
		s = ""
	} else {
		s = parseMultiValue(output)
	}
	return strings.ReplaceAll(expr, "$1", s)
}

func parseMultiValue(output string) string {
	builder := strings.Builder{}
	sOutout := strings.Split(output, ",")
	for index, val := range sOutout {
		if index != 0 {
			if index > 0 {
				builder.WriteString(",")
			}
		}
		if len(val) > 0 {
			builder.WriteString("'" + val + "'")
		}
	}
	return builder.String()
}

var exprSanitizePermission ExprSanitize = func(output, expr string) string {
	return SanitizeRegExOutPut(output, expr)
}

func SanitizeRegExOutPut(output, expr string) string {
	if strings.Contains(output, common.GrepRegex) {
		output = ""
	}
	return  strings.ReplaceAll(expr, "$1", output)
 }
