package models

import (
	"encoding/json"
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
	Name       string      `json:"name"`
	AuditTests []AuditTest `json:"audit_tests"`
}

//AuditTest data model
type AuditTest struct {
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
func (at *AuditTest) UnmarshalJSON(data []byte) error {
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
	default:
		at.Sanitize = exprSanitizePermission
	}
	return nil
}

//ExprSanitize sanitize expr
type ExprSanitize func(expr string) string

var exprSanitizeOwnership ExprSanitize = func(expr string) string {
	return validateRegExOutPut(expr)
}

var exprSanitizeProcessParam ExprSanitize = func(expr string) string {
	return validateRegExOutPut(expr)
}

var exprSanitizePermission ExprSanitize = func(expr string) string {
	return validateRegExOutPut(expr)
}

func validateRegExOutPut(expr string) string {
	if strings.Contains(expr, "[^\"]\\S*'") {
		return ""
	}
	return expr
}
