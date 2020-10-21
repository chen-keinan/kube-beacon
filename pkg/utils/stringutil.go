package utils

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"strconv"
	"strings"
)

//ExprSanitize sanitize expr
type ExprSanitize func(output []string, expr string) string

//ExprSanitizeMultiProcessParam check type
var ExprSanitizeMultiProcessParam ExprSanitize = func(outputArr []string, expr string) string {
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
			value = exp.EvaFunc(output, i, exp.Expr)
			exp.Expr = value
		}
		builder.WriteString(value)
	}
	return builder.String()
}

var parseMultiValue EvalFunction = func(output string, index int, expr string) string {
	//add condition value before split to array
	variable := fmt.Sprintf("'$%s'", strconv.Itoa(index))
	if strings.Contains(expr, variable) {
		fOutPut := fmt.Sprintf("'%s'", output)
		return strings.ReplaceAll(expr, variable, fOutPut)
	}
	sOutout := strings.Split(output, ",")
	if len(sOutout) == 1 {
		return sanitizeSingleValue(expr, index, sOutout[0])
	}
	return sanitizeMultiValue(sOutout, index, expr)
}

func sanitizeMultiValue(sOutout []string, index int, expr string) string {
	builderOne := strings.Builder{}
	for index, val := range sOutout {
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

func sanitizeSingleValue(expr string, index int, sOutout string) string {
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

//parseSingleValue for regex case
var parseSingleValue EvalFunction = func(output string, index int, expr string) string {
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
			exprList = append(exprList, Expr{Type: common.MultiValue, Expr: s, EvaFunc: parseMultiValue})
		} else {
			exprList = append(exprList, Expr{Type: common.SingleValue, Expr: s, EvaFunc: parseSingleValue})
		}
	}
	return exprList
}

//EvalFunction evaluate command result with expression
type EvalFunction func(output string, index int, expr string) string

//Expr data
type Expr struct {
	Type    string
	Expr    string
	EvaFunc EvalFunction
}
