package templateutil

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

type emptyInventory struct{}

type inventory struct {
	Material string
	Count    uint
}

func TestEvaluateTemplateStringToString(t *testing.T) {
	t.Run("empty template, empty inventory", func(t *testing.T) {
		got, err := EvaluateTemplateStringToString("", emptyInventory{}, template.FuncMap{})
		require.NoError(t, err)
		require.Equal(t, "", got)
	})

	t.Run("empty template, populated inventory", func(t *testing.T) {
		got, err := EvaluateTemplateStringToString("", inventory{Material: "wool", Count: 17}, template.FuncMap{})
		require.NoError(t, err)
		require.Equal(t, "", got)
	})

	t.Run("plain string without template expressions", func(t *testing.T) {
		got, err := EvaluateTemplateStringToString("no template string", inventory{Material: "wool", Count: 17}, template.FuncMap{})
		require.NoError(t, err)
		require.Equal(t, "no template string", got)
	})

	t.Run("missing field on empty inventory returns error", func(t *testing.T) {
		_, err := EvaluateTemplateStringToString(
			"{{.Count}} items are made of {{.Material}}",
			emptyInventory{}, template.FuncMap{},
		)
		require.Error(t, err)
	})

	t.Run("simple substitution with default delimiters", func(t *testing.T) {
		got, err := EvaluateTemplateStringToString(
			"{{.Count}} items are made of {{.Material}}",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
		)
		require.NoError(t, err)
		require.Equal(t, "17 items are made of wool", got)
	})

	t.Run("custom FuncMap is available inside template", func(t *testing.T) {
		funcs := template.FuncMap{
			"isOne": func(i int) bool { return i == 1 },
		}
		got, err := EvaluateTemplateStringToString(
			"{{if isOne 1}}yes{{else}}no{{end}}",
			emptyInventory{}, funcs,
		)
		require.NoError(t, err)
		require.Equal(t, "yes", got)
	})
}

func TestEvaluateTemplateStringToStringWithDelimiter(t *testing.T) {
	t.Run("custom delimiters substitute", func(t *testing.T) {
		got, err := EvaluateTemplateStringToStringWithDelimiter(
			"<<.Count>> items are made of <<.Material>>",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
			"<<", ">>",
		)
		require.NoError(t, err)
		require.Equal(t, "17 items are made of wool", got)
	})

	t.Run("default markers are literal when custom delimiters set", func(t *testing.T) {
		got, err := EvaluateTemplateStringToStringWithDelimiter(
			"{{.Count}} items are made of {{.Material}}",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
			"<<", ">>",
		)
		require.NoError(t, err)
		require.Equal(t, "{{.Count}} items are made of {{.Material}}", got)
	})
}

func TestEvaluateTemplateStringToStringWithDelimiterAndOpts(t *testing.T) {
	t.Run("no options, custom delimiters", func(t *testing.T) {
		got, err := EvaluateTemplateStringToStringWithDelimiterAndOpts(
			"<<.Count>> items are made of <<.Material>>",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
			"<<", ">>",
			[]string{},
		)
		require.NoError(t, err)
		require.Equal(t, "17 items are made of wool", got)
	})

	t.Run("missingkey=error surfaces missing field as error", func(t *testing.T) {
		got, err := EvaluateTemplateStringToStringWithDelimiterAndOpts(
			"<<.Undefined>> items are made of <<.Material>>",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
			"<<", ">>",
			[]string{"missingkey=error"},
		)
		require.EqualError(t, err, `template: :1:2: executing "" at <.Undefined>: can't evaluate field Undefined in type templateutil.inventory`)
		require.Equal(t, "", got)
	})

	t.Run("missingkey=error also applies with default delimiters", func(t *testing.T) {
		got, err := EvaluateTemplateStringToStringWithDelimiterAndOpts(
			"{{.UndefinedDefDelim}} items are made of {{.Material}}",
			inventory{Material: "wool", Count: 17}, template.FuncMap{},
			"", "",
			[]string{"missingkey=error"},
		)
		require.EqualError(t, err, `template: :1:2: executing "" at <.UndefinedDefDelim>: can't evaluate field UndefinedDefDelim in type templateutil.inventory`)
		require.Equal(t, "", got)
	})
}

func TestEvaluateTemplateStringToString_parseError(t *testing.T) {
	_, err := EvaluateTemplateStringToString("{{ unterminated", inventory{Material: "wool", Count: 17}, template.FuncMap{})
	require.Error(t, err)
}
