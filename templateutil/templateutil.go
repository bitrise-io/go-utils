// Package templateutil provides helpers for evaluating Go text/template
// snippets against an inventory value, with optional custom delimiters
// and template options.
package templateutil

import (
	"bytes"
	"text/template"
)

// evaluateTemplate parses and executes templateContent against inventory
// using the provided FuncMap, delimiters, and options. Empty delimiters
// fall back to the text/template defaults ("{{" and "}}").
// See https://pkg.go.dev/text/template#Template.Option for the option set.
func evaluateTemplate(
	templateContent string,
	inventory any,
	funcs template.FuncMap,
	delimLeft, delimRight string,
	templateOptions []string,
) (string, error) {
	tmpl := template.New("").Funcs(funcs).Delims(delimLeft, delimRight)
	if len(templateOptions) > 0 {
		tmpl = tmpl.Option(templateOptions...)
	}

	tmpl, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, inventory); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// EvaluateTemplateStringToStringWithDelimiterAndOpts evaluates
// templateContent with custom delimiters and template options applied.
// See https://pkg.go.dev/text/template#Template.Option for the option set.
func EvaluateTemplateStringToStringWithDelimiterAndOpts(
	templateContent string,
	inventory any,
	funcs template.FuncMap,
	delimLeft, delimRight string,
	templateOptions []string,
) (string, error) {
	return evaluateTemplate(templateContent, inventory, funcs, delimLeft, delimRight, templateOptions)
}

// EvaluateTemplateStringToStringWithDelimiter evaluates templateContent
// with custom delimiters and no template options.
func EvaluateTemplateStringToStringWithDelimiter(
	templateContent string,
	inventory any,
	funcs template.FuncMap,
	delimLeft, delimRight string,
) (string, error) {
	return evaluateTemplate(templateContent, inventory, funcs, delimLeft, delimRight, nil)
}

// EvaluateTemplateStringToString evaluates templateContent using the
// default text/template delimiters and no template options.
func EvaluateTemplateStringToString(
	templateContent string,
	inventory any,
	funcs template.FuncMap,
) (string, error) {
	return evaluateTemplate(templateContent, inventory, funcs, "", "", nil)
}
