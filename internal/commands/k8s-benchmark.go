package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
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
		res, err := ls.Execute()
		if err != nil {
			fmt.Printf("Failed to execute command %s", err.Error())
		}
		outputs := strings.Split(res.Stdout, "\n")
		bk.evalExpression(outputs, at)
	}
}

func (bk K8sBenchmark) evalExpression(outputs []string, at k8s.AuditTest) {
	for _, o := range outputs {
		if len(o) == 0 {
			continue
		}
		expr := at.Sanitize(o) + at.EvalExpr
		expression, err := govaluate.NewEvaluableExpression(expr)
		if err != nil {
			fmt.Print(emoji.Sprintf("audit test %s :cross_mark:\n", at.Description))
		}
		result, err := expression.Evaluate(nil)
		if err != nil {
			fmt.Print(emoji.Sprintf("audit test %s :cross_mark:\n", at.Description))
		}
		if result.(bool) {
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
