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
	var s string
	if strings.Contains(output, common.GrepRegex) {
		if strings.Contains(expr, "'$1'") {
			expr = strings.ReplaceAll(expr, "'$1'", "$1")
		}
		s = "''"
		return strings.ReplaceAll(expr, "$1", s)
	}
	return parseMultiValue(output, expr)

}

func parseMultiValue(output, expr string) string {
	//add condition value before split to array
	if strings.Contains(expr, "'$1'") {
		expr = strings.ReplaceAll(expr, "'$1'", "'"+output+"'")
	}
	sOutout := strings.Split(output, ",")
	if len(sOutout) == 1 {
		return sanitizeSingleValue(expr, sOutout)
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

func sanitizeSingleValue(expr string, sOutout []string) string {
	if strings.Contains(expr, "IN") {
		expr = strings.ReplaceAll(expr, "IN", "==")
	}
	return strings.ReplaceAll(expr, "($1)", "'"+sOutout[0]+"'")
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
