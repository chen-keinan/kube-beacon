package benchmark

import (
	"embed"
	"fmt"
	"github.com/chen-keinan/beacon/pkg/utils"
	"io/ioutil"
)

const K8sFolder = "k8s/v1.6.0"
const GkeFolder = "gke/v1.1.0"

var (
	//go:embed k8s/v1.6.0
	resK8s embed.FS

	//go:embed gke/v1.1.0
	resGke embed.FS
)

func LoadK8sSpecs() ([]utils.FilesInfo, error) {
	dir, _ := resK8s.ReadDir(K8sFolder)
	specs := make([]utils.FilesInfo, 0)
	for _, r := range dir {
		file, err := resK8s.Open(fmt.Sprintf("%s/%s", K8sFolder, r.Name()))
		if err != nil {
			return specs, err
		}
		data, err := ioutil.ReadAll(file)
		spec := utils.FilesInfo{Name: r.Name(), Data: string(data)}
		if err != nil {
			return specs, err
		}
		if err != nil {
			return specs, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func LoadGkeSpecs() ([]utils.FilesInfo, error) {
	dir, _ := resGke.ReadDir(GkeFolder)
	specs := make([]utils.FilesInfo, 0)
	for _, r := range dir {
		file, err := resGke.Open(fmt.Sprintf("%s/%s", GkeFolder, r.Name()))
		if err != nil {
			return specs, err
		}
		data, err := ioutil.ReadAll(file)
		spec := utils.FilesInfo{Name: r.Name(), Data: string(data)}
		if err != nil {
			return specs, err
		}
		if err != nil {
			return specs, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}
