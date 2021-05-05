package main

import (
	"fmt"
	"github.com/chen-keinan/beacon/pkg/models"
)

//K8sBenchAuditResultHook this plugin method accept k8s audit bench results
//event include test data , description , audit, remediation and result
//nolint
func K8sBenchAuditResultHook(k8sAuditResults models.KubeAuditResults) error {
	fmt.Println("this is K8sBenchAuditResultHook plugin")
	return nil
}
