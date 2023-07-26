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
		"stats_initial.go",
		"stats_looped.go",
		"stats_goroutined.go",
		"stats_goroutined_fastjson.go",
		"stats_alternate.go",
		"stats_benchmark_test.go",
		"stats_common_benchmark_test.go",
		"stats_common_optimization_test.go",
		"stats_common_test.go",
		"stats_common_test.out",
		"stats_common_benchmark_test.out",
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
