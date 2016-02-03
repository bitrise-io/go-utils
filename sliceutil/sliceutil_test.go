package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
