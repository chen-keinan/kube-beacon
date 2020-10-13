package utils

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_CheckType_Permission_OK test
func Test_CheckType_Permission_OK(t *testing.T) {
	evalExpr := "$1 <= 644"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("700", 1, evalExpr)
	assert.Equal(t, ti, "700 <= 644")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Owner_OK(t *testing.T) {
	evalExpr := "'$1' == 'root:root';"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("root:root", 1, evalExpr)
	assert.Equal(t, ti, "'root:root' == 'root:root'")
}

//Test_CheckType_ProcessParam_OK test
func Test_CheckType_ProcessParam_OK(t *testing.T) {
	evalExpr := "'$1' == 'false';"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("false", 1, evalExpr)
	assert.Equal(t, ti, "'false' == 'false'")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_OK(t *testing.T) {
	evalExpr := "'RBAC' IN ($1);"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("RBAC,bbb", 1, evalExpr)
	assert.Equal(t, ti, "'RBAC' IN ('RBAC','bbb')")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_RexOK(t *testing.T) {
	evalExpr := "'RBAC' IN ($1);"
	bench := ExprSanitizeMultiProcessParam
	ti := bench(common.GrepRegex, 1, evalExpr)
	assert.Equal(t, ti, "'RBAC' == ''")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Regex_OK(t *testing.T) {
	evalExpr := "'$1' == 'root:root';"
	bench := ExprSanitizeMultiProcessParam
	ti := bench(common.GrepRegex, 1, evalExpr)
	assert.Equal(t, ti, "'' == 'root:root'")
}

//Test_CheckType_Regex_MultiParamType test
func Test_CheckType_Regex_MultiParamType(t *testing.T) {
	evalExpr := "'$1' != 'root:root'; && 'root:root' IN ($1);"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("root:root", 1, evalExpr)
	assert.Equal(t, ti, "'root:root' != 'root:root' && 'root:root' == 'root:root'")
}

//Test_CheckType_Regex_MultiParamTypeManyValues test
func Test_CheckType_Regex_MultiParamTypeManyValues(t *testing.T) {
	evalExpr := "'$1' != 'root:root'; && 'root:root' IN ($1);"
	bench := ExprSanitizeMultiProcessParam
	ti := bench("root:root,abc", 1, evalExpr)
	assert.Equal(t, ti, "'root:root,abc' != 'root:root' && 'root:root' IN ('root:root','abc')")
}
