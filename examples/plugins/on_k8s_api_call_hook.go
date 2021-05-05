package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/beacon/pkg/models"
)

//K8sBenchAuditResultHook this plugin method accept k8s audit bench results
//event include test data , description , audit, remediation and result
func K8sBenchAuditResultHook(k8sAuditResults models.KubeAuditResults) error {
	var sb = new(bytes.Buffer)
	err := json.NewEncoder(sb).Encode(k8sAuditResults)
	fmt.Print(k8sAuditResults)
	if err != nil {
		return err
	}
	return nil
}
