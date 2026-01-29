package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFirstLine(t *testing.T) {
	t.Log("Empty input")
	{
		require.Equal(t, "", ReadFirstLine("", false))
		require.Equal(t, "", ReadFirstLine("", true))
	}

	t.Log("Multiline empty input - ignore-empty-lines:false")
	{
		firstLine := ReadFirstLine(`


`, false)
		require.Equal(t, "", firstLine)
	}

	t.Log("Multiline empty input - ignore-empty-lines:true")
	{
		firstLine := ReadFirstLine(`


`, true)
		require.Equal(t, "", firstLine)
	}

	t.Log("Multiline non empty input - ignore-empty-lines:false")
	{
		firstLine := ReadFirstLine(`first line

second line`, false)
		require.Equal(t, "first line", firstLine)
	}

	t.Log("Multiline empty input - ignore-empty-lines:true")
	{
		firstLine := ReadFirstLine(`first line

second line`, true)
		require.Equal(t, "first line", firstLine)
	}

	t.Log("Multiline non empty input, with leading empty line - ignore-empty-lines:false")
	{
		firstLine := ReadFirstLine(`

first line

second line`, false)
		require.Equal(t, "", firstLine)
	}

	t.Log("Multiline non empty input, with leading empty line - ignore-empty-lines:true")
	{
		firstLine := ReadFirstLine(`

first line

second line`, true)
		require.Equal(t, "first line", firstLine)
	}
}

func TestGenericTrim(t *testing.T) {
	require.Equal(t, "", genericTrim("", 4, false, false))
	require.Equal(t, "", genericTrim("", 4, false, true))
	require.Equal(t, "", genericTrim("", 4, true, false))
	require.Equal(t, "", genericTrim("", 4, true, true))

	require.Equal(t, "1234", genericTrim("123456789", 4, false, false))
	require.Equal(t, "1...", genericTrim("123456789", 4, false, true))
	require.Equal(t, "6789", genericTrim("123456789", 4, true, false))
	require.Equal(t, "...9", genericTrim("123456789", 4, true, true))
}

func TestMaxLastCharsWithDots(t *testing.T) {
	require.Equal(t, "", MaxLastCharsWithDots("", 10))
	require.Equal(t, "", MaxLastCharsWithDots("1234", 1))
	require.Equal(t, "...56", MaxLastCharsWithDots("123456", 5))
	require.Equal(t, "123456", MaxLastCharsWithDots("123456", 6))
	require.Equal(t, "123456", MaxLastCharsWithDots("123456", 10))

	require.Equal(t, "... world!", MaxLastCharsWithDots("hello world!", 10))
}

func TestMaxFirstCharsWithDots(t *testing.T) {
	require.Equal(t, "", MaxFirstCharsWithDots("", 10))
	require.Equal(t, "", MaxFirstCharsWithDots("1234", 1))
	require.Equal(t, "12...", MaxFirstCharsWithDots("123456", 5))
	require.Equal(t, "123456", MaxFirstCharsWithDots("123456", 6))
	require.Equal(t, "123456", MaxFirstCharsWithDots("123456", 10))

	require.Equal(t, "hello w...", MaxFirstCharsWithDots("hello world!", 10))
}

func TestLastNLines(t *testing.T) {
	t.Log("Empty input")
	{
		require.Equal(t, "", LastNLines("", 1))
		require.Equal(t, "", LastNLines("", 5))
	}

	t.Log("Single line")
	{
		require.Equal(t, "line", LastNLines("line", 1))
		require.Equal(t, "line", LastNLines("line", 5))
	}

	t.Log("Multiple lines")
	{
		input := `line1
line2
line3
line4
line5`
		require.Equal(t, "line5", LastNLines(input, 1))
		require.Equal(t, "line4\nline5", LastNLines(input, 2))
		require.Equal(t, "line3\nline4\nline5", LastNLines(input, 3))
		require.Equal(t, input, LastNLines(input, 5))
		require.Equal(t, input, LastNLines(input, 10))
	}

	t.Log("Lines with surrounding newlines")
	{
		input := `
line1
line2
`
		require.Equal(t, "line2", LastNLines(input, 1))
		require.Equal(t, "line1\nline2", LastNLines(input, 2))
		require.Equal(t, "line1\nline2", LastNLines(input, 3))
	}

	t.Log("N = 0")
	{
		input := `line1
line2`
		require.Equal(t, "", LastNLines(input, 0))
	}
}

func TestIndentTextWithMaxLength(t *testing.T) {
	t.Log("Empty")
	{
		input := ""
		output := IndentTextWithMaxLength(input, "", 80, true)
		require.Equal(t, "", output)
	}

	t.Log("One liner")
	{
		input := "one liner"
		output := IndentTextWithMaxLength(input, "", 80, true)
		require.Equal(t, "one liner", output)
	}

	t.Log("One liner - with indent")
	{
		input := "one liner"
		output := IndentTextWithMaxLength(input, " => ", 76, true)
		require.Equal(t, " => one liner", output)
	}

	t.Log("One liner - max width")
	{
		input := "one"
		output := IndentTextWithMaxLength(input, "", 3, true)
		require.Equal(t, "one", output)
	}

	t.Log("One liner - longer than max width")
	{
		input := "onetwo"
		output := IndentTextWithMaxLength(input, "", 3, true)
		require.Equal(t, "one\ntwo", output)
	}

	t.Log("One liner - max width - with indent")
	{
		input := "one"
		require.Equal(t, " on\n e", IndentTextWithMaxLength(input, " ", 2, true))
		require.Equal(t, "on\n e", IndentTextWithMaxLength(input, " ", 2, false))
	}

	t.Log("One liner - max width - with first-line indent false")
	{
		input := "one"
		output := IndentTextWithMaxLength(input, " ", 2, false)
		require.Equal(t, "on\n e", output)
	}

	t.Log("One liner - longer than max width - with indent")
	{
		input := "onetwo"
		require.Equal(t, " on\n et\n wo", IndentTextWithMaxLength(input, " ", 2, true))
		require.Equal(t, "on\n et\n wo", IndentTextWithMaxLength(input, " ", 2, false))
	}

	t.Log("Two lines, shorter than max")
	{
		input := `first line
second line`
		output := IndentTextWithMaxLength(input, "", 80, true)
		require.Equal(t, `first line
second line`, output)
	}

	t.Log("Two lines, shorter than max - with indent")
	{
		input := `first line
second line`

		require.Equal(t, `  first line
  second line`, IndentTextWithMaxLength(input, "  ", 78, true))
		require.Equal(t, `first line
  second line`, IndentTextWithMaxLength(input, "  ", 78, false))
	}

	t.Log("Two lines, longer than max")
	{
		input := `firstline
secondline`
		output := IndentTextWithMaxLength(input, "", 5, true)
		require.Equal(t, `first
line
secon
dline`, output)
	}

	t.Log("Max length = 0")
	{
		input := "Indent is longer than max length"
		require.Equal(t, "", IndentTextWithMaxLength(input, "...", 0, true))
		require.Equal(t, "", IndentTextWithMaxLength(input, "...", 0, false))
	}

	t.Log("Max length < 0")
	{
		input := "text"
		require.Equal(t, "", IndentTextWithMaxLength(input, "", -1, true))
	}
}
