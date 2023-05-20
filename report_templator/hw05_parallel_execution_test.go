package reporttemplator

import (
	"fmt"
	"testing"
)

// go run -test templator.go hw05_parallel_execution.go

func TestHw05(t *testing.T) {
	_ = t
	directory := "../hw05_parallel_execution/"
	files := []string{
		"run.go",
		"run_test.go",
		"statistic.go",
		"statistic_test.go",
		"TestRun4TaskWith5Gorutine.txt",
		"TestRunAllTasksWithoutAnyError.txt",
		"TestRunFirstMTasksErrors.txt",
		"TestRunWithUnlimitedErrorsCount.txt",
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
