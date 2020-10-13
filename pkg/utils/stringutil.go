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
	sExpr := separateExpr(expr)
	for _, exp := range sExpr {
		for i, output := range outputArr {
			if !strings.Contains(exp.Expr, "$") {
				continue
			}
			if exp.Type == common.SingleValue {
				value = sanitizeRegExOutPut(output, i+1, exp.Expr)
			} else {
				value = parseMultiValue(output, i+1, exp.Expr)
			}
			exp.Expr = value
		}
		builder.WriteString(value)
	}
	return builder.String()
}

func parseMultiValue(output string, index int, expr string) string {
	//add condition value before split to array
	variable := fmt.Sprintf("'$%s'", strconv.Itoa(index))
	if strings.Contains(expr, variable) {
		fOutPut := fmt.Sprintf("'%s'", output)
		expr = strings.ReplaceAll(expr, variable, fOutPut)
	}
	sOutout := strings.Split(output, ",")
	if len(sOutout) == 1 {
		return sanitizeSingleValue(expr, index, sOutout[0])
	}
	return sanitizeMultiValue(sOutout, expr)
}

func sanitizeMultiValue(sOutout []string, expr string) string {
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
	return strings.ReplaceAll(expr, "$1", builderOne.String())
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

//sanitizeRegExOutPut for regex case
func sanitizeRegExOutPut(output string, index int, expr string) string {
	if strings.Contains(output, common.GrepRegex) {
		output = ""
	}
	varaible := fmt.Sprintf("$%s", strconv.Itoa(index))
	return strings.ReplaceAll(expr, varaible, output)
}

func separateExpr(expr string) []Expr {
	exprList := make([]Expr, 0)
	split := strings.Split(expr, ";")
	for _, s := range split {
		if len(s) == 0 {
			continue
		}
		if strings.Contains(s, "IN") {
			exprList = append(exprList, Expr{Type: common.MultiValue, Expr: s})
		} else {
			exprList = append(exprList, Expr{Type: common.SingleValue, Expr: s})
		}
	}
	return exprList
}

//Expr data
type Expr struct {
	Type string
	Expr string
}
