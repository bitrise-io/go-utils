package templateutil

import (
	"testing"

	"text/template"

	"github.com/stretchr/testify/require"
)

type EmptyInventory struct {
}

type Inventory struct {
	Material string
	Count    uint
}

func TestIndexOfStringInSlice(t *testing.T) {
	t.Log("Empty")
	result, err := EvaluateTemplateStringToString("", EmptyInventory{}, template.FuncMap{})
	require.NoError(t, err)
	require.Equal(t, "", result)
	//
	result, err = EvaluateTemplateStringToString("", Inventory{"wool", 17}, template.FuncMap{})
	require.NoError(t, err)
	require.Equal(t, "", result)
	//
	result, err = EvaluateTemplateStringToString("no template string", Inventory{"wool", 17}, template.FuncMap{})
	require.NoError(t, err)
	require.Equal(t, "no template string", result)

	t.Log("Empty inventory - missing argument/property (error)")
	result, err = EvaluateTemplateStringToString("{{.Count}} items are made of {{.Material}}",
		EmptyInventory{}, template.FuncMap{})
	require.Error(t, err)

	//
	var templateFuncMap = template.FuncMap{
		"isOne": func(i int) bool {
			return i == 1
		},
	}

	t.Log("Simple")
	inv := Inventory{"wool", 17}
	result, err = EvaluateTemplateStringToString("{{.Count}} items are made of {{.Material}}",
		inv, template.FuncMap{})
	require.NoError(t, err)
	require.Equal(t, "17 items are made of wool", result)

	inv = Inventory{"glass", 18}
	result, err = EvaluateTemplateStringToString("{{.Count}} items are made of {{.Material}}",
		inv, templateFuncMap)
	require.NoError(t, err)
	require.Equal(t, "18 items are made of glass", result)
}
