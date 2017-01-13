package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCommandSlice(t *testing.T) {
	t.Log("it fails if slice empty")
	{
		cmd, err := NewFromSlice()
		require.Error(t, err)
		require.Equal(t, (*Model)(nil), cmd)
	}

	t.Log("it creates cmd if cmdSlice has 1 element")
	{
		_, err := NewFromSlice("ls")
		require.NoError(t, err)
	}

	t.Log("it creates cmd if cmdSlice has multiple elements")
	{
		_, err := NewFromSlice("ls", "-a", "-l", "-h")
		require.NoError(t, err)
	}
}
