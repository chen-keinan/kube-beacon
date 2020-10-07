package commands

import (
	"encoding/json"
	"fmt"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
	"strconv"
	"strings"
)

//K8sBenchmark k8s benchmark object
type K8sBenchmark struct {
}

//Help return benchmark command help
func (bk K8sBenchmark) Help() string {
	return "-a , --audit run benchmark audit tests"
}

//Run execute benchmark command
func (bk K8sBenchmark) Run(args []string) int {
	audit := k8s.Audit{}
	auditFile := utils.GetK8sBenchmarkAuditTestsFile()
	err := json.Unmarshal([]byte(auditFile), &audit)
	if err != nil {
		fmt.Print("Failed to read audit test file")
	}
	for _, ac := range audit.Categories {
		bk.runTests(ac)
	}

	return 0
}

func (bk K8sBenchmark) runTests(ac k8s.Category) {
	for _, at := range ac.SubCategory.AuditTests {
		ls := execute.ExecTask{
			Command: at.AuditCommand,
			Args:    []string{},
			Shell:   true,
		}
		fmt.Println(at.AuditCommand)
		res, err := ls.Execute()
		if err != nil {
			fmt.Printf("Failed to execute command %s", err.Error())
		}
		outputs := strings.Split(res.Stdout, "\n")
		switch at.CheckType {
		case "permission":
			bk.checkPermission(outputs, res, at)
		case "ownership":
			bk.checkOwnership(outputs, at)
		}
	}
}

func (bk K8sBenchmark) checkPermission(outputs []string, res execute.ExecResult, at k8s.AuditTest) {
	for _, o := range outputs {
		if len(o) == 0 {
			continue
		}
		value, err := strconv.Atoi(o)
		if err != nil {
			fmt.Println(res.Stderr)
		}
		if value <= 644 {
			fmt.Print(emoji.Sprintf("audit test %s :check_mark_button:\n", at.Description))
		} else {
			fmt.Print(emoji.Sprintf("audit test %s :cross_mark:\n", at.Description))
		}
	}
}

func (bk K8sBenchmark) checkOwnership(outputs []string, at k8s.AuditTest) {
	for _, o := range outputs {
		if len(o) == 0 {
			continue
		}
		if o == "root:root" {
			fmt.Print(emoji.Sprintf("audit test %s :check_mark_button:\n", at.Description))
		} else {
			fmt.Print(emoji.Sprintf("audit test %s :cross_mark:\n", at.Description))
		}
	}
}

//Synopsis for help
func (bk K8sBenchmark) Synopsis() string {
	return bk.Help()
}
