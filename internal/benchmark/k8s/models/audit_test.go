package models

import (
	"encoding/json"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	CheckTypePermission        = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"Master Node Configuration Files\",\"audit_tests\":[{\"name\":\"Ensure that the API server pod specification file permissions are set to 644 or more restrictive (Automated)\",\"description\":\"Ensure that the API server pod specification file has permissions of 644 or more restrictive.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml\",\"remediation\":\"chmod 644 /etc/kubernetes/manifests/kube-apiserver.yaml\",\"check_type\":\"permission\",\"impact\":\"None\",\"eval_expr\":\"$1 <= 644\",\"default_value\":\"By default, the kube-apiserver.yaml file has permissions of 640.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\"]}]}}]}"
	CheckTypeOwner             = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"Master Node Configuration Files\",\"audit_tests\":[{\"name\":\"Ensure that the API server pod specification file ownership is set to root:root (Automated)\",\"description\":\"Ensure that the API server pod specification file ownership is set to root:root.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"stat -c %U:%G /etc/kubernetes/manifests/kube-apiserver.yaml\",\"remediation\":\"chown root:root /etc/kubernetes/manifests/kube-apiserver.yaml\",\"check_type\":\"ownership\",\"impact\":\"None\",\"eval_expr\":\"'$1' == 'root:root'\",\"default_value\":\"By default, the kube-apiserver.yaml file ownership is set to root:root.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\"]}]}}]}"
	CheckTypeProcessParam      = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --anonymous-auth argument is set to false (Manual)\",\"description\":\"Disable anonymous requests to the API server.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"ps -ef | grep kube-apiserver |grep 'anonymous-auth' | grep -o 'anonymous-auth=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\",\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the below parameter.\\n--anonymous-auth=false\",\"check_type\":\"process_param\",\"impact\":\"Anonymous requests will be rejected.\",\"eval_expr\":\"'$1' == 'false'\",\"default_value\":\"By default, anonymous access is enabled.\",\"references\":[\"https://kubernetes.io/docs/admin/kube-apiserver/\",\"https://kubernetes.io/docs/admin/authentication/#anonymous-requests\"]}]}}]}"
	CheckTypeMultiProcessParam = "{\"benchmark_type\":\"k8s\",\"categories\":[{\"name\":\"Control Plane Components\",\"sub_category\":{\"name\":\"API Server\",\"audit_tests\":[{\"name\":\"Ensure that the --authorization-mode argument includes RBAC (Automated)\",\"description\":\"Turn on Role Based Access Control.\",\"profile_applicability\":\"Level 1 - Master Node\",\"audit\":\"ps -ef | grep kube-apiserver |grep 'authorization-mode' | grep -o 'authorization-mode=[^\\\"]\\\\S*' | awk -F \\\"=\\\" '{print $2}' |awk 'FNR <= 1'\",\"remediation\":\"Edit the API server pod specification file /etc/kubernetes/manifests/kube- apiserver.yaml on the master node and set the --authorization-mode parameter to a value that includes RBAC, for example:--authorization-mode=Node,RBAC\",\"check_type\":\"multi_process_param\",\"impact\":\"When RBAC is enabled you will need to ensure that appropriate RBAC settings (including Roles, RoleBindings and ClusterRoleBindings) are configured to allow appropriate access.\",\"eval_expr\":\"'RBAC' IN ($1)\",\"default_value\":\"By default, RBAC authorization is not enabled.\",\"references\":[\"https://kubernetes.io/docs/reference/access-authn-authz/rbac/\"]}]}}]}"
)

//Test_CheckType_Permission_OK test
func Test_CheckType_Permission_OK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypePermission), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize("700", bench.EvalExpr)
	assert.Equal(t, ti, "700 <= 644")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Owner_OK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeOwner), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize("root:root", bench.EvalExpr)
	assert.Equal(t, ti, "'root:root' == 'root:root'")
}

//Test_CheckType_ProcessParam_OK test
func Test_CheckType_ProcessParam_OK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize("false", bench.EvalExpr)
	assert.Equal(t, ti, "'false' == 'false'")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_OK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize("RBAC,bbb", bench.EvalExpr)
	assert.Equal(t, ti, "'RBAC' IN ('RBAC','bbb')")
}

//Test_CheckType_Multi_ProcessParam_OK test
func Test_CheckType_Multi_ProcessParam_RexOK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeMultiProcessParam), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize(common.GrepRegex, bench.EvalExpr)
	assert.Equal(t, ti, "'RBAC' IN ('')")
}

//Test_CheckType_Owner_OK test
func Test_CheckType_Regex_OK(t *testing.T) {
	ab := Audit{}
	err := json.Unmarshal([]byte(CheckTypeOwner), &ab)
	if err != nil {
		t.Fatal(err)
	}
	bench := ab.Categories[0].SubCategory.AuditTests[0]
	ti := bench.Sanitize(common.GrepRegex, bench.EvalExpr)
	assert.Equal(t, ti, "'' == 'root:root'")
}

//Test_JsonMarshalBad test
func Test_JsonMarshalBad(t *testing.T) {
	ab := &AuditBench{}
	err := ab.UnmarshalJSON([]byte("{bad,wwq}"))
	if err == nil {
		t.Fatal(err)
	}
}
