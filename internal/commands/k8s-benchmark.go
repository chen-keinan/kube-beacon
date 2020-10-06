package commands

import (
	"encoding/json"
	"fmt"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s"
	"github.com/kyokomi/emoji"
	"strconv"
)

const auditJson = `{
  "benchmark_type": "k8s",
  "categories": [
    {
      "name": "Control Plane Components",
      "sub_category": {
        "name": "Master Node Configuration Files",
        "audit_tests": [
          {
            "name": "Ensure that the API server pod specification file permissions are set to 644 or more restrictive (Automated)",
            "description": "Ensure that the API server pod specification file has permissions of 644 or more restrictive.",
            "profile_applicability": "Level 1 - Master Node",
            "audit": "stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml",
            "remediation": "chmod 644 /etc/kubernetes/manifests/kube-apiserver.yaml",
            "check_type": "permission"
          },
          {
            "name": "Ensure that the API server pod specification file ownership is set to root:root (Automated)",
            "description": "Ensure that the API server pod specification file ownership is set to root:root.",
            "profile_applicability": "Level 1 - Master Node",
            "audit": "stat -c %U:%G /etc/kubernetes/manifests/kube-apiserver.yaml",
            "remediation": "chown root:root /etc/kubernetes/manifests/kube-apiserver.yaml",
            "check_type": "ownership"
          }
        ]
      }
    }
  ]
}`

type K8sBenchmark struct {
}

func (*K8sBenchmark) Help() string {
	return "-a , --audit run benchmark audit tests"
}
func (*K8sBenchmark) Run(args []string) int {
	audit := k8s.Audit{}
	err := json.Unmarshal([]byte(auditJson), &audit)
	if err != nil {
		fmt.Print("Failed to read audit test file")
	}
	for _, ac := range audit.Categories {
		for _, at := range ac.SubCategory.AuditTests {
			ls := execute.ExecTask{
				Command: at.AuditCommand,
				Args:    []string{},
				Shell:   true,
			}
			res, err := ls.Execute()

			if err != nil {
				fmt.Sprintf("Failed to execute command %s", err.Error())
			}
 			switch at.CheckType {
			case "permission":
				value, err := strconv.Atoi(res.Stdout)
				if err != nil {
					fmt.Println(res.Stderr)
				}
				if value <= 644 {
					fmt.Print(emoji.Sprintf("audit test %s :check_mark_button:\n", at.Description))
				}
			case "ownership":
				fmt.Println(res.Stderr)
				if res.Stdout == "root:root" {
					fmt.Print(emoji.Sprintf("audit test %s :check_mark_button:\n", at.Description))
				}else{
					fmt.Print(emoji.Sprintf("audit test %s :check_mark_button:\n", at.Description))
				}
			}
		}
	}

	return 0
}
func (h *K8sBenchmark) Synopsis() string {
	return h.Help()
}
