package reporttemplator

import (
	"fmt"
	"testing"
)

// go test templator.go hw11_telnet_client_test.go

func TestHw11(t *testing.T) {
	_ = t
	directory := "../hw11_telnet_client/"
	files := []string{
		"test.native.sh",
		"test.native.sh.out",
		"test.sh.out",
		"telnet_test.go.txt",
		"main_test.go.txt",
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
