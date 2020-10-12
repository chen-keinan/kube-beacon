package commands

import (
	"encoding/json"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

const CheckTypeMultiProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"ps -ef | grep kube-apiserver |grep 'authorization-mode' | grep -o 'authorization-mode=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\",\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_process_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"!('RBAC' IN ($1))\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
const CheckTypeMultiExprProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"1.2.11 Ensure that the admission control plugin AlwaysAdmit is not set\",\"description\":\"Do not allow all requests.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"ps -ef | grep kube-apiserver |grep 'enable-admission-plugins' | grep -o 'enable-admission-plugins=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\",\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and either remove the --enable-admission-plugins parameter, or set it to a value that does not include AlwaysAdmit.\",\"check_type\":\"multi_process_param\",\"impact\":\"Only requests explicitly allowed by the admissions control plugins would be served.\",\"eval_expr\":\"'$1' != '' && !('AlwaysAdmit' IN ($1))\",\"default_value\":\"AlwaysAdmit is not in the list of default admission plugins.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\",\"https://kubernetes.io/docs/admin/admission-controllers/#alwaysadmit\"]}]}}]}"

//Test_EvalVarSingleIn text
func Test_EvalVarSingleIn(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, err := kb.evalExpression([]string{"aaa"}, bench)
	assert.True(t, match)
	assert.NoError(t, err)
}

//Test_EvalVarSingleNotInGood text
func Test_EvalVarSingleNotInGood(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, mErr := kb.evalExpression([]string{"ttt,aaa"}, bench)
	assert.True(t, match)
	assert.NoError(t, mErr)
}

//Test_EvalVarSingleNotInBad text
func Test_EvalVarSingleNotInBad(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, mErr := kb.evalExpression([]string{"RBAC,aaa"}, bench)
	assert.False(t, match)
	assert.NoError(t, mErr)
}

//Test_EvalVarSingleNotInSingleValue test
func Test_EvalVarSingleNotInSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, mErr := kb.evalExpression([]string{"aaa"}, bench)
	assert.True(t, match)
	assert.NoError(t, mErr)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprSingleValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiExprProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, mErr := kb.evalExpression([]string{"AlwaysAdmit"}, bench)
	assert.False(t, match)
	assert.NoError(t, mErr)
}

//Test_EvalVarMultiExprSingleValue test
func Test_EvalVarMultiExprMultiValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiExprProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	match, mErr := kb.evalExpression([]string{"bbb,aaa"}, bench)
	assert.True(t, match)
	assert.NoError(t, mErr)
}
