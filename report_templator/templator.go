package main

import (
	"os"
	"strings"
)

type Template struct {
	content       string
	substitutions map[string]string
}

func (template *Template) loadFromFile(filepath string, with_escaping bool) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		template.content = ""
	}
	template.content = string(data)
	if with_escaping {
		template.content = escaping(template.content)
	}
	return err
}

func (template *Template) render() string {
	result := template.content
	for k := range template.substitutions {
		result = strings.Replace(result, "{{ "+k+" }}", template.substitutions[k], -1)
	}
	return result
}

func escaping(content string) string {
	content = tab_escaping(content)
	content = strings.Replace(content, "\n```\n", "\n'''\n", -1)
	content = strings.Replace(content, "\n```text\n", "\n'''text\n", -1)
	content = strings.Replace(content, "\n```go\n", "\n'''go\n", -1)
	content = strings.Replace(content, "{{ ", "{"+string('\x02')+"{ ", -1)
	content = strings.Replace(content, " }}", " }"+string('\x02')+"}", -1)
	return content
}

func tab_escaping(content string) string {
	content = strings.Replace(content, "\t", "    ", -1)
	return content
}

func (template *Template) renderToFile(filepath string) error {
	f, errCreate := os.Create(filepath)
	if errCreate != nil {
		return errCreate
	}
	result := template.render()
	defer f.Close()

	_, errWrite := f.WriteString(result)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
