package bplugin

import (
	m2 "github.com/chen-keinan/beacon/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func Test_NewPluginWorker(t *testing.T) {
	production, err := zap.NewProduction()
	assert.NoError(t, err)
	completedChan := make(chan bool)
	plChan := make(chan m2.KubeAuditResults)
	pw := NewPluginWorker(NewPluginWorkerData(plChan, K8sBenchAuditResultHook{}, completedChan), production)
	assert.True(t, len(pw.cmd.plugins.Plugins) == 0)

}
