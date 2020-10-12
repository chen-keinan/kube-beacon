package utils

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_CheckType_Permission_OK test
func Test_CheckType_Permission_OK(t *testing.T) {
	evalExpr := "$1 <= 644"
	bench := ExprSanitizePermission
	ti := bench("700", evalExpr)
	assert.Equal(t, ti, "700 <= 644")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Owner_OK(t *testing.T) {
	evalExpr := "'$1' == 'root:root'"
	bench := ExprSanitizeOwnership
	ti := bench("root:root", evalExpr)
	assert.Equal(t, ti, "'root:root' == 'root:root'")
}

//Test_CheckType_ProcessParam_OK test
func Test_CheckType_ProcessParam_OK(t *testing.T) {
	evalExpr := "'$1' == 'false'"
	bench := ExprSanitizeProcessParam
	ti := bench("false", evalExpr)
	assert.Equal(t, ti, "'false' == 'false'")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_OK(t *testing.T) {
	evalExpr := "'RBAC' IN ($1)"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("RBAC,bbb", evalExpr)
	assert.Equal(t, ti, "'RBAC' IN ('RBAC','bbb')")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_RexOK(t *testing.T) {
	evalExpr := "'RBAC' IN ($1)"
	bench := ExprSanitizeMultiProcessParam
	ti := bench(common.GrepRegex, evalExpr)
	assert.Equal(t, ti, "'RBAC' IN ('')")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Regex_OK(t *testing.T) {
	evalExpr := "'$1' == 'root:root'"
	bench := ExprSanitizeOwnership
	ti := bench(common.GrepRegex, evalExpr)
	assert.Equal(t, ti, "'' == 'root:root'")
}
