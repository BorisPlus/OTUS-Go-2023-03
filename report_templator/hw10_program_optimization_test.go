package reporttemplator

import (
	"fmt"
	"testing"
)

// go test templator.go hw10_program_optimization_test.go

func TestHw10(t *testing.T) {
	_ = t
	directory := "../hw10_program_optimization/"
	files := []string{
		"stats_example.go",
		"stat.go",
		"experimantal/main.go",
		"stats_remark.go",
		"stats_benchmark_test.go",
	}

	substitutions := make(map[string]string)
	for _, file := range files {
		tmplt := Template{}
		tmplt.loadFromFile(fmt.Sprintf("%s%s", directory, file))
		substitutions[file] = tmplt.render(true)
	}

	report := Template{}
	report.loadFromFile(fmt.Sprintf("%s%s", directory, "REPORT.template.md"))
	report.substitutions = substitutions
	report.renderToFile(fmt.Sprintf("%s%s", directory, "REPORT.md"))
}
