package filters

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

//Test_SpecificTestPredict text
func Test_SpecificTestPredict(t *testing.T) {
	ab := []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}
	abp := SpecificTest(ab, "1.2.1")
	assert.Equal(t, abp[0].Name, "1.2.1 abc")
	assert.True(t, len(abp) == 1)
}

//Test_SpecificTestPredicateNoValidArg text
func Test_SpecificTestPredicateNoValidArg(t *testing.T) {
	ab := []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}
	abp := SpecificTest(ab, "1.2.5")
	assert.Equal(t, abp[0].Name, "1.2.1 abc")
	assert.Equal(t, abp[1].Name, "1.2.3 eft")
	assert.True(t, len(abp) == 2)
}

//Test_BasicPredicate text
func Test_BasicPredicate(t *testing.T) {
	ab := []*models.AuditBench{{Name: "1.2.1 abc"}, {Name: "1.2.3 eft"}}
	abp := Basic(ab, "")
	assert.Equal(t, abp[0].Name, "1.2.1 abc")
	assert.Equal(t, abp[1].Name, "1.2.3 eft")
	assert.True(t, len(abp) == 2)
}
