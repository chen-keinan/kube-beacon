package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/benchmark/k8s/models"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
	"strings"
)

//K8sAudit k8s benchmark object
type K8sAudit struct {
}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return "-a , --audit run benchmark audit tests"
}

//Run execute benchmark command
func (bk K8sAudit) Run(args []string) int {
	audit := models.Audit{}
	auditFiles := utils.GetK8sBenchmarkAuditTestsFile()
	for _, auditFile := range auditFiles {
		err := json.Unmarshal([]byte(auditFile), &audit)
		if err != nil {
			fmt.Print("Failed to read audit test file")
		}
		for _, ac := range audit.Categories {
			bk.runTests(ac)
		}
	}
	return 0
}

func (bk K8sAudit) runTests(ac models.Category) {
	for _, at := range ac.SubCategory.AuditTests {
		result, err := shell.NewShellExec().Exec(at.AuditCommand)
		if err != nil {
			fmt.Printf("Failed to execute command %s", err.Error())
			continue
		}
		outputs := strings.Split(result.Stdout, "\n")
		bk.evalExpression(outputs, at)
	}
}

func (bk K8sAudit) evalExpression(outputs []string, at models.AuditTest) {
	for _, o := range outputs {
		if len(o) == 0 && len(outputs) > 1 {
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
func (bk K8sAudit) Synopsis() string {
	return bk.Help()
}
