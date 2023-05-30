package reporttemplator

import (
	"fmt"
	"testing"
)

// go run -test templator.go hw06_pipeline_execution.go

func TestHw06(t *testing.T) {
	_ = t
	directory := "../hw06_pipeline_execution/"
	files := []string{
		"pipeline.go",
		"N100TimesTesting.txt",
		"TestPipelineConcurencyTime.txt",
		"TestNoStages.txt",
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
