package reporttemplator

import (
	"os"
	"strings"
)

type Template struct {
	content       string
	substitutions map[string]string
}

func (template *Template) loadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		template.content = ""
	}
	template.content = string(data)
	return err
}

func (template *Template) render(withEscaping bool) string {
	result := template.content
	if withEscaping {
		result = tabEscaping(result)
	}
	for k := range template.substitutions {
		result = strings.ReplaceAll(result, "{{ "+k+" }}", template.substitutions[k])
	}
	return result
}

// func escaping(content string) string {
// 	content = tabEscaping(content)
// 	content = strings.ReplaceAll(content, "\n```\n", "\n'''\n")
// 	content = strings.ReplaceAll(content, "\n```text\n", "\n'''text\n")
// 	content = strings.ReplaceAll(content, "\n```go\n", "\n'''go\n")
// 	content = strings.ReplaceAll(content, "{{ ", "{"+string('\x02')+"{ ")
// 	content = strings.ReplaceAll(content, " }}", " }"+string('\x02')+"}")
// 	return content
// }

func tabEscaping(content string) string {
	content = strings.ReplaceAll(content, "\t", "    ")
	return content
}

func (template *Template) renderToFile(filepath string) error {
	f, errCreate := os.Create(filepath)
	if errCreate != nil {
		return errCreate
	}
	result := template.render(false)
	defer f.Close()

	_, errWrite := f.WriteString(result)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
