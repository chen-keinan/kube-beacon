package utils

import (
	"github.com/chen-keinan/beacon/internal/common"
	"strings"
)

//ExprSanitize sanitize expr
type ExprSanitize func(output, expr string) string

//ExprSanitizeOwnership check type
var ExprSanitizeOwnership ExprSanitize = func(output, expr string) string {
	return sanitizeRegExOutPut(output, expr)
}

//ExprSanitizeProcessParam check type
var ExprSanitizeProcessParam ExprSanitize = func(output, expr string) string {
	return sanitizeRegExOutPut(output, expr)
}

//ExprSanitizeMultiProcessParam check type
var ExprSanitizeMultiProcessParam ExprSanitize = func(output, expr string) string {
	var value string
	builder := strings.Builder{}
	sExpr := separateExpr(expr)
	for _, exp := range sExpr {
		if exp.Type == common.SingleValue {
			value = sanitizeRegExOutPut(output, exp.Expr)
		} else {
			value = parseMultiValue(output, exp.Expr)
		}
		builder.WriteString(value)
	}
	return builder.String()
}

func parseMultiValue(output, expr string) string {
	//add condition value before split to array
	if strings.Contains(expr, "'$1'") {
		expr = strings.ReplaceAll(expr, "'$1'", "'"+output+"'")
	}
	sOutout := strings.Split(output, ",")
	if len(sOutout) == 1 {
		return sanitizeSingleValue(expr, sOutout[0])
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

func sanitizeSingleValue(expr string, sOutout string) string {
	if strings.Contains(expr, "IN") {
		expr = strings.ReplaceAll(expr, "IN", "==")
	}
	if sOutout == common.GrepRegex {
		sOutout = ""
	}
	return strings.ReplaceAll(expr, "($1)", "'"+sOutout+"'")
}

//ExprSanitizePermission check type
var ExprSanitizePermission ExprSanitize = func(output, expr string) string {
	return sanitizeRegExOutPut(output, expr)
}

//sanitizeRegExOutPut for regex case
func sanitizeRegExOutPut(output, expr string) string {
	if strings.Contains(output, common.GrepRegex) {
		output = ""
	}
	return strings.ReplaceAll(expr, "$1", output)
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
