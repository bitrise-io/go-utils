package command_test

import (
	"os/exec"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/require"
)

func TestNewCommandSlice(t *testing.T) {
	t.Log("it fails if slice empty")
	{
		cmd, err := command.NewFromSlice([]string{})
		require.Error(t, err)
		require.Equal(t, (*command.Model)(nil), cmd)
	}

	t.Log("it creates cmd if cmdSlice has 1 element")
	{
		_, err := command.NewFromSlice([]string{"ls"})
		require.NoError(t, err)
	}

	t.Log("it creates cmd if cmdSlice has multiple elements")
	{
		_, err := command.NewFromSlice([]string{"ls", "-a", "-l", "-h"})
		require.NoError(t, err)
	}
}

func TestNewWithParams(t *testing.T) {
	t.Log("it fails if params empty")
	{
		cmd, err := command.NewWithParams()
		require.Error(t, err)
		require.Equal(t, (*command.Model)(nil), cmd)
	}

	t.Log("it creates cmd if params has 1 element")
	{
		_, err := command.NewWithParams("ls")
		require.NoError(t, err)
	}

	t.Log("it creates cmd if params has multiple elements")
	{
		_, err := command.NewWithParams("ls", "-a", "-l", "-h")
		require.NoError(t, err)
	}
}
