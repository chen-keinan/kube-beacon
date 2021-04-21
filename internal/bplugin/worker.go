package bplugin

import (
	"fmt"
	"github.com/chen-keinan/beacon/pkg/models"
	"go.uber.org/zap"
)

//PluginWorker instance which match command data to specific pattern
type PluginWorker struct {
	cmd *PluginWorkerData
	log *zap.Logger
}

//NewPluginWorker return new plugin worker instance
func NewPluginWorker(commandMatchData *PluginWorkerData, log *zap.Logger) *PluginWorker {
	return &PluginWorker{cmd: commandMatchData, log: log}
}

//NewPluginWorkerData return new plugin worker instance
func NewPluginWorkerData(plChan chan models.KubeAuditResults, hook K8sBenchAuditResultHook) *PluginWorkerData {
	return &PluginWorkerData{plChan: plChan, plugins: hook}
}

//PluginWorkerData encapsulate plugin worker properties
type PluginWorkerData struct {
	plChan  chan models.KubeAuditResults
	plugins K8sBenchAuditResultHook
}

//Invoke invoke plugin accept audit bench results
func (pm *PluginWorker) Invoke() {
	go func() {
		ae := <-pm.cmd.plChan
		if len(pm.cmd.plugins.Plugins) > 0 {
			for _, pl := range pm.cmd.plugins.Plugins {
				err := ExecuteK8sAuditResults(pl, ae)
				if err != nil {
					pm.log.Error(fmt.Sprintf("failed to execute plugins %s", err.Error()))
				}
			}
		}
	}()
}
