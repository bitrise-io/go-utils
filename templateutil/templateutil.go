package templateutil

import (
	"bytes"
	"text/template"
)

// EvaluateTemplateStringToStringWithDelimiter ...
func EvaluateTemplateStringToStringWithDelimiter(templateContent string, inventory interface{}, funcs template.FuncMap, delimLeft, delimRight string) (string, error) {
	tmpl := template.New("").Funcs(funcs).Delims(delimLeft, delimRight)
	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", err
	}

	var resBuffer bytes.Buffer
	if err := tmpl.Execute(&resBuffer, inventory); err != nil {
		return "", err
	}

	return resBuffer.String(), nil
}

// EvaluateTemplateStringToString ...
func EvaluateTemplateStringToString(templateContent string, inventory interface{}, funcs template.FuncMap) (string, error) {
	return EvaluateTemplateStringToStringWithDelimiter(templateContent, inventory, funcs, "", "")
}
