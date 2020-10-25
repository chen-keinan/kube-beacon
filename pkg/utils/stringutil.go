package utils

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"strconv"
	"strings"
)

//CmdExprBuilder sanitize expr
type CmdExprBuilder func(output []string, expr string) string

//UpdateCmdExprParam check type
var UpdateCmdExprParam CmdExprBuilder = func(outputArr []string, expr string) string {
	var value string
	builder := strings.Builder{}
	sExpr := SeparateExpr(expr)
	for _, exp := range sExpr {
		for i, output := range outputArr {
			if !strings.Contains(exp.Expr, "$") {
				if i > 0 {
					break
				} else {
					value = exp.Expr
					break
				}
			}
			value = exp.EvaExprBuilderFunc(strings.TrimSpace(output), i, exp.Expr)
			exp.Expr = value
		}
		builder.WriteString(value)
	}
	return builder.String()
}

//parseMultiValue build evaluation expresion for expr with IN clause
var parseMultiValue EvaExprBuilderFunc = func(output string, index int, expr string) string {
	//add condition value before split to array
	variable := fmt.Sprintf("'$%s'", strconv.Itoa(index))
	if strings.Contains(expr, variable) {
		fOutPut := fmt.Sprintf("'%s'", output)
		return strings.ReplaceAll(expr, variable, fOutPut)
	}
	sOutput := strings.Split(output, ",")
	if len(sOutput) == 1 {
		return changeExprFromMultiToSingle(expr, index, sOutput[0])
	}
	return buildInClauseExpr(sOutput, index, expr)
}

func buildInClauseExpr(sOutPut []string, index int, expr string) string {
	builderOne := strings.Builder{}
	for index, val := range sOutPut {
		if index != 0 {
			if index > 0 {
				builderOne.WriteString(",")
			}
		}
		if len(val) > 0 {
			builderOne.WriteString("'" + val + "'")
		}
	}
	return strings.ReplaceAll(expr, fmt.Sprintf("$%d", index), builderOne.String())
}

//changeExprFromMultiToSingle it change the expression from multi to single
// where IN clause has only one param
func changeExprFromMultiToSingle(expr string, index int, sOutout string) string {
	variable := fmt.Sprintf("($%s)", strconv.Itoa(index))
	fOutPut := fmt.Sprintf("'%s'", sOutout)
	if strings.Contains(expr, "IN") {
		expr = strings.ReplaceAll(expr, "IN", "==")
	}
	if sOutout == common.GrepRegex {
		fOutPut = "''"
	}
	return strings.ReplaceAll(expr, variable, fOutPut)
}

//parseSingleValue build evaluation expresion for expr for non IN clause
var parseSingleValue EvaExprBuilderFunc = func(output string, index int, expr string) string {
	if strings.Contains(output, common.GrepRegex) {
		output = ""
	}
	varaible := fmt.Sprintf("$%s", strconv.Itoa(index))
	return strings.ReplaceAll(expr, varaible, output)
}

//SeparateExpr separate expression to single and multi blocks
func SeparateExpr(expr string) []Expr {
	exprList := make([]Expr, 0)
	split := strings.Split(expr, ";")
	for _, s := range split {
		if len(s) == 0 {
			continue
		}
		if strings.Contains(s, "IN") && strings.Contains(s, "$") {
			exprList = append(exprList, Expr{Type: common.MultiValue, Expr: s, EvaExprBuilderFunc: parseMultiValue})
		} else {
			exprList = append(exprList, Expr{Type: common.SingleValue, Expr: s, EvaExprBuilderFunc: parseSingleValue})
		}
	}
	return exprList
}

//EvaExprBuilderFunc build evaluation expresion
//it replace expression params with audit command result
type EvaExprBuilderFunc func(output string, index int, expr string) string

//Expr data
type Expr struct {
	Type               string
	Expr               string
	EvaExprBuilderFunc EvaExprBuilderFunc
}

//ExcludeAuditTest return true if test is not included in specific tests to run
func ExcludeAuditTest(tests []string, name string) bool {
	if len(tests) == 0 {
		return false
	}
	for _, t := range tests {
		if strings.Contains(name, t) {
			return false
		}
	}
	return true
}

//GetSpecificTestsToExecute return processing function by specificTests
func GetSpecificTestsToExecute(arg string) []string {
	values := strings.ReplaceAll(arg, "s=", "")
	return strings.Split(values, ",")
}

//RemoveNewLineSuffix remove new line from suffix
func RemoveNewLineSuffix(str string) string {
	i := len(str)
	if len(str) > 0 && str[i-1:i] == "\n" {
		return str[0 : i-1]
	}
	return str
}

//AddNewLineToNonEmptyStr add new line to non empty string
func AddNewLineToNonEmptyStr(str string) string {
	if !strings.HasSuffix(str, "\n") {
		return fmt.Sprintf("%s\n", str)
	}
	return str
}
