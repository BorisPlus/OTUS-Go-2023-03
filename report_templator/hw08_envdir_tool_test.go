package reporttemplator

import (
	"fmt"
	"testing"
)

// go test templator.go hw08_envdir_tool_test.go

func TestHw08(t *testing.T) {
	_ = t
	directory := "../hw08_envdir_tool/"
	files := []string{
		"hw08_go_doc_-all.txt",
		"hw08_go_test_-v.txt",
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
