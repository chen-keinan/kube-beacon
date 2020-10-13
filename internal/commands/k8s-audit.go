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
	auditFiles, err := utils.GetK8sBenchAuditFiles()
	if err != nil {
		panic(fmt.Sprintf("failed to read audit files %s", err))
	}
	for _, auditFile := range auditFiles {
		err := json.Unmarshal([]byte(auditFile.Data), &audit)
		if err != nil {
			panic("Failed to unmarshal audit test json file")
		}
		for _, ac := range audit.Categories {
			bk.runTests(ac)
		}
	}
	return 0
}

func (bk K8sAudit) runTests(ac models.Category) {
	for _, at := range ac.SubCategory.AuditTests {
		resArr := make([]string, 0)
		for _, val := range at.AuditCommand {
			result, err := shell.NewShellExec().Exec(val)
			if err != nil {
				fmt.Printf("Failed to execute command %s", err.Error())
				continue
			}
			resArr = append(resArr, result.Stdout)
		}
		combResult := make([]string, 0)
		MapResult := make(map[string]int)
		match, err := bk.evalExpression(resArr, len(resArr), combResult, 0, at, len(resArr), MapResult)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			bk.printTestResults(match, at)
		}
	}
}

func (bk K8sAudit) printTestResults(match bool, at models.AuditBench) {
	if match {
		fmt.Print(emoji.Sprintf(":check_mark_button: %s\n", at.Name))
	} else {
		fmt.Print(emoji.Sprintf(":cross_mark: %s\n", at.Name))

	}
}

func (bk K8sAudit) evalExpression(resArr []string, index int, combArr []string, count int, at models.AuditBench, origSize int, resMap map[string]int) (bool, error) {
	bk.ValidMultiExpr(resArr, index, combArr, count, at, origSize, resMap)
	return resMap["match"] == resMap["total"], nil
}

//ValidMultiExpr validate multi expr recursive
func (bk K8sAudit) ValidMultiExpr(resArr []string, index int, combArr []string, count int, at models.AuditBench, origSize int, resMap map[string]int) {
	if len(resArr) == 0 {
		return
	}
	outputs := strings.Split(resArr[count], "\n")
	for _, o := range outputs {
		if len(o) == 0 && len(outputs) > 1 {
			continue
		}
		combArr = append(combArr, o)
		bk.ValidMultiExpr(resArr[count+1:index], index-1, combArr, count, at, origSize, resMap)
		if origSize == len(combArr) {
			expr := at.Sanitize(combArr, at.EvalExpr)
			resMap["total"]++
			count, err := bk.evalCommandExpr(at, expr)
			if err != nil {
				fmt.Println(err)
			}
			resMap["match"] += count
		}
		combArr = combArr[:len(combArr)-1]
	}

}

func (bk K8sAudit) evalCommandExpr(at models.AuditBench, expr string) (int, error) {

	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to build evaluation command expr for\n %s", at.Name)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate command expr for audit test %s", at.Name)
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}

//Synopsis for help
func (bk K8sAudit) Synopsis() string {
	return bk.Help()
}
