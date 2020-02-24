package sliceutil

import (
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/stretchr/testify/require"
)

func TestUniqueStringSlice(t *testing.T) {
	require.Equal(t, []string{}, UniqueStringSlice([]string{}))
	require.Equal(t, []string{"one"}, UniqueStringSlice([]string{"one"}))
	testutil.EqualSlicesWithoutOrder(t,
		[]string{"one", "two"},
		UniqueStringSlice([]string{"one", "two"}))
	testutil.EqualSlicesWithoutOrder(t,
		[]string{"one", "two", "three"},
		UniqueStringSlice([]string{"one", "two", "three", "two", "one"}))
}

func TestIndexOfStringInSlice(t *testing.T) {
	t.Log("Empty slice")
	require.Equal(t, -1, IndexOfStringInSlice("abc", []string{}))

	testSlice := []string{"abc", "def", "123", "456", "123"}

	t.Log("Find item")
	require.Equal(t, 0, IndexOfStringInSlice("abc", testSlice))
	require.Equal(t, 1, IndexOfStringInSlice("def", testSlice))
	require.Equal(t, 3, IndexOfStringInSlice("456", testSlice))

	t.Log("Find first item, if multiple")
	require.Equal(t, 2, IndexOfStringInSlice("123", testSlice))

	t.Log("Item is not in the slice")
	require.Equal(t, -1, IndexOfStringInSlice("cba", testSlice))
}

func TestIsStringInSlice(t *testing.T) {
	t.Log("Empty slice")
	require.Equal(t, false, IsStringInSlice("abc", []string{}))

	testSlice := []string{"abc", "def", "123", "456", "123"}

	t.Log("Find item")
	require.Equal(t, true, IsStringInSlice("abc", testSlice))
	require.Equal(t, true, IsStringInSlice("def", testSlice))
	require.Equal(t, true, IsStringInSlice("456", testSlice))

	t.Log("Find first item, if multiple")
	require.Equal(t, true, IsStringInSlice("123", testSlice))

	t.Log("Item is not in the slice")
	require.Equal(t, false, IsStringInSlice("cba", testSlice))
}

func TestCleanWhitespace(t *testing.T) {
	// Arrange
	tests := []struct {
		name      string
		arg       []string
		omitEmpty bool
		want      []string
	}{
		{
			name:      "empty list",
			arg:       []string(nil),
			omitEmpty: true,
			want:      []string(nil),
		},
		{
			name:      "multiple elements",
			arg:       []string{"url1", "url2", "url3"},
			omitEmpty: true,
			want:      []string{"url1", "url2", "url3"},
		},
		{
			name:      "multiple elements, skipping empty ones",
			arg:       []string{"url1", "url2", " \n ", "url3"},
			omitEmpty: true,
			want:      []string{"url1", "url2", "url3"},
		},
		{
			name: "multiple elements with spaces and newlines",
			arg: []string{"url1", `
url2   `, `

url3`},
			omitEmpty: true,
			want:      []string{"url1", "url2", "url3"},
		},
		{
			name:      "multiple elements not omitting empty elements",
			arg:       []string{"url1", " \n ", "url3"},
			omitEmpty: false,
			want:      []string{"url1", "", "url3"},
		},
	}

	// Act + Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotItems := CleanWhitespace(tt.arg, tt.omitEmpty); !reflect.DeepEqual(gotItems, tt.want) {
				t.Errorf("%v, want %v", gotItems, tt.want)
			}
		})
	}
}
