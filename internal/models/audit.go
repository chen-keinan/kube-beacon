package models

import (
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"strings"
)

//Audit data model
type Audit struct {
	BenchmarkType string     `yaml:"benchmark_type"`
	Categories    []Category `yaml:"categories"`
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
	CmdExprBuilder       utils.CmdExprBuilder
	TestSucceed          bool
	CommandParams        map[int][]string
	Category             string
	NonApplicable        bool `mapstructure:"non_applicable" yaml:"non_applicable"`
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
	switch at.CheckType {
	case "multi_param":
		at.CmdExprBuilder = utils.UpdateCmdExprParam
	}
	at.CommandParams = make(map[int][]string)
	for index, command := range at.AuditCommand {
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
