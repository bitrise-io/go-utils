package cmdex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCommandSlice(t *testing.T) {
	t.Log("it fails if cmdSlice empty")
	{
		cmd, err := NewCommandFromSlice([]string{})
		require.Error(t, err)
		require.Equal(t, (*CommandModel)(nil), cmd)
	}

	t.Log("it creates cmd if cmdSlice has 1 element")
	{
		_, err := NewCommandFromSlice([]string{"ls"})
		require.NoError(t, err)
	}

	t.Log("it creates cmd if cmdSlice has multiple elements")
	{
		_, err := NewCommandFromSlice([]string{"ls", "-a", "-l", "-h"})
		require.NoError(t, err)
	}
}
