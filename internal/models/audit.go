package models

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/mitchellh/mapstructure"
)

//Audit data model
type Audit struct {
	BenchmarkType string     `yaml:"benchmark_type"`
	Categories    []Category `yaml:"categories"`
}

//AuditTestTotals model
type AuditTestTotals struct {
	Warn int
	Pass int
	Fail int
}

//Category data model
type Category struct {
	Name        string       `yaml:"name"`
	SubCategory *SubCategory `yaml:"sub_category"`
}

//SubCategory data model
type SubCategory struct {
	Name       string        `yaml:"name"`
	AuditTests []*AuditBench `yaml:"audit_tests"`
}

//AuditBench data model
type AuditBench struct {
	Name                 string   `mapstructure:"name" yaml:"name"`
	ProfileApplicability string   `mapstructure:"profile_applicability" yaml:"profile_applicability"`
	Description          string   `mapstructure:"description" yaml:"description"`
	AuditCommand         []string `mapstructure:"audit" json:"audit"`
	CheckType            string   `mapstructure:"check_type" yaml:"check_type"`
	Remediation          string   `mapstructure:"remediation" yaml:"remediation"`
	Impact               string   `mapstructure:"impact" yaml:"impact"`
	DefaultValue         string   `mapstructure:"default_value" yaml:"default_value"`
	References           []string `mapstructure:"references" yaml:"references"`
	EvalExpr             string   `mapstructure:"eval_expr" yaml:"eval_expr"`
	TestSucceed          bool
	CommandParams        map[int][]string
	Category             string
	NonApplicable        bool
	TestType             string `mapstructure:"type" yaml:"type"`
}

//AuditResult data
type AuditResult struct {
	NumOfExec    int
	NumOfSuccess int
}

//UnmarshalYAML over unmarshall to add logic
func (at *AuditBench) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var res map[string]interface{}
	if err := unmarshal(&res); err != nil {
		return err
	}
	err := mapstructure.Decode(res, &at)
	if err != nil {
		return err
	}
	at.CommandParams = make(map[int][]string)
	if at.TestType == common.NonApplicableTest || at.TestType == common.ManualTest {
		at.NonApplicable = true
	}
	return nil
}
