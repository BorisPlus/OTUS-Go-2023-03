package reporttemplator

import (
	"fmt"
	"testing"
)

// go test templator.go hw09_struct_validator_test.go

func TestHw09(t *testing.T) {
	_ = t
	directory := "../hw09_struct_validator/"
	files := []string{
		"go_doc_-all.txt",
		"TestValidatePositive.txt",
		"TestValidateNegative.txt",
		"TestValidateNotStructObject.txt",
		"TestValidateNotImplemented.txt",
		"TestValidateExpectedNotImplemented.txt",
		"TestValidateUnxpectedNotImplemented.txt",
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
