package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/beacon/pkg/models"
	"net/http"
	"strings"
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
	req, err := http.NewRequest("POST", "http://localhost:8090/audir-results", strings.NewReader(sb.String()))
	if err != nil {
		return err
	}
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}
