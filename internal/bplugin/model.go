package bplugin

import "plugin"

//K8sBenchAuditResultHook hold the plugin symbol for K8s bench audit result Hook
type K8sBenchAuditResultHook struct {
	Plugins []plugin.Symbol
}
