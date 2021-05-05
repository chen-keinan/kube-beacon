package bplugin

import (
	"fmt"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/pkg/models"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestPluginLoader_Plugins(t *testing.T) {
	pl, err := pluginSetUp("k8s_bench_audit_result_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	assert.Equal(t, plFiles[0], "test_plugin.go")
}


func TestExecuteNetEvt(t *testing.T) {
	pl, err := pluginSetUp("k8s_bench_audit_result_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	sym, err := pl.Compile(plFiles[0], common.K8sBenchAuditResultHook)
	assert.NoError(t, err)
		err = ExecuteK8sAuditResults(sym, models.KubeAuditResults{})
	assert.NoError(t, err)
}

func TestPluginLoader_CompileBad(t *testing.T) {
	pl, err := pluginSetUp("empty.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	_, err = pl.Compile(plFiles[0], common.K8sBenchAuditResultHook)
	assert.Error(t, err)
	_, err = pl.Compile("a/b/c", common.K8sBenchAuditResultHook)
	assert.Error(t, err)
}
func TestPluginLoader_CompileWrongHook(t *testing.T) {
	pl, err := pluginSetUp("k8s_bench_audit_result_hook.go")
	assert.NoError(t, err)
	plFiles, err := pl.Plugins()
	assert.NoError(t, err)
	_, err = pl.Compile(plFiles[0], "NoHook")
	assert.Error(t, err)
}

func pluginSetUp(fileName string) (*PluginLoader, error) {
	fm := utils.NewKFolder()
	folder, err := utils.GetPluginSourceSubFolder(fm)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(folder)
	if err != nil {
		return nil, err
	}
	cfolder, err := utils.GetCompilePluginSubFolder(fm)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(cfolder)
	if err != nil {
		return nil, err
	}
	err = utils.CreateHomeFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	err = utils.CreatePluginsSourceFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	err = utils.CreatePluginsCompiledFolderIfNotExist(fm)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("./fixtures/%s", fileName))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	nf, err := os.Create(path.Join(folder, "test_plugin.go"))
	if err != nil {
		return nil, err
	}
	_, err = nf.WriteString(string(data))
	if err != nil {
		return nil, err
	}
	pl, err := NewPluginLoader()
	if err != nil {
		return nil, err
	}
	return pl, err
}
