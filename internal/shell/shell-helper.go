package shell

import (
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var log = logger.GetLog()

// #nosec
func isProcessOwnerAdmin() bool {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(stdout)) == common.RootUser
}

//UpdateProcessOwnerConfig copy kubectl config to admin folder if the owner is admin
func UpdateProcessOwnerConfig() {
	if isProcessOwnerAdmin() {
		_, err := NewShellExec().Exec("cp /etc/kubernetes/admin.conf /root/.kube/config")
		if err != nil {
			log.Console("failed to copy kubectl config to root folder")
		}
	}
}
