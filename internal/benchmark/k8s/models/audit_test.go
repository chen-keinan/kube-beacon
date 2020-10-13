package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	CheckTypeMultiProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'authorization-mode' | grep -o 'authorization-mode=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"'RBAC' IN ($1)\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
	CheckTypeBlaa              = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":[\"ps -ef | grep kube-apiserver |grep 'authorization-mode' | grep -o 'authorization-mode=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\"],\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"blaa\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"'RBAC' IN ($1)\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
)

//Test_CheckType_Permission test
func Test_CheckType_Blaa(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeBlaa), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Nil(t, ab.Categories[0].SubCategory.AuditTests[0].Sanitize)
}

//Test_CheckType_Multi_ProcessParam test
func Test_CheckType_Multi_ProcessParam(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, ab.Categories[0].SubCategory.AuditTests[0].Sanitize)
}
