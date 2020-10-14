package commands

import (
	"encoding/json"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s/models"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/mocks"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

const CheckTypeMultiProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'authorization-mode' | grep -o 'authorization-mode=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"!('RBAC' IN ($0));\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
const CheckTypeMultiExprProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"1.2.11 Ensure that the admission control plugin AlwaysAdmit is not set\",\"description\":\"Do not allow all requests.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'enable-admission-plugins' | grep -o 'enable-admission-plugins=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and either remove the --enable-admission-plugins parameter, or set it to a value that does not include AlwaysAdmit.\",\"check_type\":\"multi_param\",\"impact\":\"Only requests explicitly allowed by the admissions control plugins would be served.\",\"eval_expr\":\"'$0' != ''; && !('AlwaysAdmit' IN ($0));\",\"default_value\":\"AlwaysAdmit is not in the list of default admission plugins.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\",\"https://kubernetes.io/docs/admin/admission-controllers/#alwaysadmit\"]}]}}]}"
const CheckTypeMultiExprEmptyProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"1.2.14 Ensure that the admission control plugin ServiceAccount is set\",\"description\":\"Automate service accounts management.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'disable-admission-plugins' | grep -o 'disable-admission-plugins=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Follow the documentation and create ServiceAccount objects as per your environment. Then, edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and ensure that the --disable-admission-plugins parameter is set to a value that does not include ServiceAccount.\",\"check_type\":\"multi_param\",\"impact\":\"None\",\"eval_expr\":\"'$0' != ''; && !('ServiceAccount' IN ($0));\",\"default_value\":\"By default, ServiceAccount is set.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\",\"https://kubernetes.io/docs/admin/admission-controllers/#serviceaccount\",\"https://kubernetes.io/docs/tasks/configure-pod-container/configure-service- account/\"]}]}}]}"
const CheckTypeComparator = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"1.2.20 Ensure that the --secure-port argument is not set to 0\",\"description\":\"Do not disable the secure port.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'secure-port' | grep -o 'secure-port=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and either remove the --secure-port parameter or set it to a different (non-zero) desired port.\\n\",\"check_type\":\"multi_param\",\"impact\":\"You need to set the API Server up with the right TLS certificates.\",\"eval_expr\":\"$0 > 0; && $0 < 65535;\",\"default_value\":\"By default, port 6443 is used as the secure port.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\"]}]}}]}"
const CheckMultiParamOK = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"aaa\",\"bbb #0\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"'$0' == '$1';\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
const CheckMultiParamNOK = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"aaa\",\"bbb #0\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"'$0' != '$1';\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"

//Test_EvalVarSingleIn text
func Test_EvalVarSingleIn(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
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
	kb.evalExpression(NewValidExprData([]string{"ttt,aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
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
	kb.evalExpression(NewValidExprData([]string{"RBAC,aaa"}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
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
	kb.evalExpression(NewValidExprData([]string{"aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
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
	kb.evalExpression(NewValidExprData([]string{"AlwaysAdmit"}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
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
	kb.evalExpression(NewValidExprData([]string{"bbb,aaa"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarMultiExprMultiEmptyValue test
func Test_EvalVarMultiExprMultiEmptyValue(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiExprEmptyProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{common.GrepRegex}, bench), make([]string, 0))
	assert.False(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_EvalVarMultiExprMultiEmptyValue test
func Test_EvalVarComparator(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckTypeComparator), &ab)
	if err != nil {
		t.Fatal(err)
	}
	kb := K8sAudit{}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	kb.evalExpression(NewValidExprData([]string{"1204"}, bench), make([]string, 0))
	assert.True(t, bench.TestResult.NumOfSuccess == bench.TestResult.NumOfExec)
	assert.NoError(t, err)
}

//Test_MultiCommandParams test
func Test_MultiCommandParams_OK(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckMultiParamOK), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.True(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}

//Test_EvalVarMultiExprMultiEmptyValue test
func Test_MultiCommandParams_NOK(t *testing.T) {
	ab := models.Audit{}
	err := json.Unmarshal([]byte(CheckMultiParamNOK), &ab)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	executor := mocks.NewMockExecutor(ctrl)
	executor.EXPECT().Exec("aaa").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	executor.EXPECT().Exec("bbb kkk").Return(&shell.CommandResult{Stdout: "kkk"}, nil).Times(1)
	kb := K8sAudit{Command: executor}
	kb.runTests(ab.Categories[0])
	assert.NoError(t, err)
	assert.False(t, ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfExec == ab.Categories[0].SubCategory.AuditTests[0].TestResult.NumOfSuccess)
}
