package main

import (
	"encoding/json"
	"fmt"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s"
	"github.com/kyokomi/emoji"
	"github.com/mitchellh/cli"
	"os"
	"strconv"
	"strings"
)

const auditJson = `{
  "benchmark_type": "k8s",
  "categories": [{
    "name": "Control Plane Components",
    "sub_category": {
      "name": "Master Node Configuration Files",
      "audit_tests": [
        {
          "name": "Ensure that the API server pod specification file permissions are set to 644 or more restrictive (Automated)",
          "description":"Ensure that the API server pod specification file has permissions of 644 or more restrictive.",
          "profile_applicability": "Level 1 - Master Node",
          "audit": "stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml",
          "remediation": "chmod 644 /etc/kubernetes/manifests/kube-apiserver.yaml"
        }
      ]
    }
  }]
}`

func main() {
	app := cli.NewCLI("beacon", "1.0.0")
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"--audit": func() (cli.Command, error) {
			return &Benchmark{}, nil
		},
		"a": func() (cli.Command, error) {
			return &Benchmark{}, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}

type Benchmark struct {
}

func (*Benchmark) Help() string {
	return "-a , --audit run benchmark tests"
}
func (*Benchmark) Run(args []string) int {
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
			value, err := strconv.Atoi(strings.Replace(res.Stderr, "stdout: ", "", -1))
			if err != nil {
				fmt.Println(res.Stderr)
				//fmt.Print("failed to convert string %s",err.Error())
			}
			if value <= 644 {
				fmt.Print(emoji.Sprintf("executing audit test %s :OK_hand:\n", at.Description))
			}
		}
	}

	return 0
}
func (h *Benchmark) Synopsis() string {
	return h.Help()
}
