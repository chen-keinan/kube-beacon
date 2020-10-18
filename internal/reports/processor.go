package reports

import (
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/jedib0t/go-pretty/table"
	"os"
)

//GenerateAuditReport generate failed audit report
func GenerateAuditReport(adts []*models.AuditBench) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Name", "Impact", "Remediation"})
	var rowCounter int
	for _, failedAudit := range adts {
		rowCounter++
		t.AppendRow(table.Row{rowCounter, failedAudit, failedAudit.Impact, failedAudit.Remediation})
		t.Render()
	}
}
