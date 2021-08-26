package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_GetSpecificTestsToExecute test
func Test_GetSpecificTestsToExecute(t *testing.T) {
	l := GetAuditTestsList("i", "i=1.2.3,1.4.5")
	assert.Equal(t, l[0], "1.2.3")
	assert.Equal(t, l[1], "1.4.5")
	l = GetAuditTestsList("e", "")
	assert.Equal(t, l[0], "")
}
