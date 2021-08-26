package utils

import (
	"fmt"
	"strings"
)

//GetAuditTestsList return processing function by specificTests
func GetAuditTestsList(key, arg string) []string {
	values := strings.ReplaceAll(arg, fmt.Sprintf("%s=", key), "")
	return strings.Split(strings.ToLower(values), ",")
}
