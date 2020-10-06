package startup

import (
	"fmt"
	"github.com/gobuffalo/packr"
)
//LoadK8sBenchmark use packr to load benchmark audit test json
func LoadK8sBenchmarkFile() string {
	box := packr.NewBox("./../benchmark/k8s")
	s, err := box.FindString("audit.json")
	if err != nil {
		panic(fmt.Sprintf("faild to load k8s benchmarks audit tests %s", err.Error()))
	}
	return s
}
