package regexputil

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedFindStringSubmatch(t *testing.T) {
	t.Log("Both the name and age group are required")
	rexp := regexp.MustCompile(`(?P<name>[a-zA-Z]+) (?P<age>[0-9]+)`)

	t.Log("Simple name+age example")
	{
		results, err := NamedFindStringSubmatch(rexp, "MyName 42")
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"name": "MyName",
			"age":  "42",
		}, results)
	}

	t.Log("Includes an additional name at the end")
	{
		results, err := NamedFindStringSubmatch(rexp, "MyName 42 AnotherName")
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"name": "MyName",
			"age":  "42",
		}, results)
	}

	t.Log("Includes an additional name at the start")
	{
		results, err := NamedFindStringSubmatch(rexp, "AnotherName MyName 42")
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"name": "MyName",
			"age":  "42",
		}, results)
	}

	t.Log("Missing name group - should error")
	{
		_, err := NamedFindStringSubmatch(rexp, " 42")
		require.EqualError(t, err, "No match found")
	}

	t.Log("Missing age group - should error")
	{
		_, err := NamedFindStringSubmatch(rexp, "MyName ")
		require.EqualError(t, err, "No match found")
	}

	t.Log("Missing both groups - should error")
	{
		_, err := NamedFindStringSubmatch(rexp, "")
		require.EqualError(t, err, "No match found")
	}

	t.Log("Optional name part")
	rexp = regexp.MustCompile(`(?P<name>[a-zA-Z]*) (?P<age>[0-9]+)`)

	t.Log("Name can now be empty - but should be included in the result!")
	{
		results, err := NamedFindStringSubmatch(rexp, " 42")
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"name": "",
			"age":  "42",
		}, results)
	}
}
