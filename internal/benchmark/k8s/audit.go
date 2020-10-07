package k8s

type Audit struct {
	BenchmarkType string     `json:"benchmark_type"`
	Categories    []Category `json:"categories"`
}

type Category struct {
	Name        string      `json:"name"`
	SubCategory SubCategory `json:"sub_category"`
}

type SubCategory struct {
	Name       string      `json:"name"`
	AuditTests []AuditTest `json:"audit_tests"`
}

type AuditTest struct {
	Name                 string   `json:"Name"`
	ProfileApplicability string   `json:"profile_applicability"`
	Description          string   `json:"description"`
	AuditCommand         string   `json:"audit"`
	CheckType            string   `json:"check_type"`
	Remediation          string   `json:"remediation"`
	Impact               string   `json:"impact"`
	DefaultValue         string   `json:"default_value"`
	References           []string `json:"references"`
}
