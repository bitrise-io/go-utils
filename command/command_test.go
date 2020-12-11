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

func TestRunCmdAndReturnExitCode(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
	}
	tests := []struct {
		name         string
		args         args
		wantExitCode int
		wantErr      bool
	}{
		{
			name: "invalid command",
			args: args{
				cmd: exec.Command(""),
			},
			wantExitCode: -1,
			wantErr:      true,
		},
		{
			name: "env command",
			args: args{
				cmd: exec.Command("env"),
			},
			wantExitCode: 0,
			wantErr:      false,
		},
		{
			name: "not existing executable",
			args: args{
				cmd: exec.Command("bash", "testdata/not_existing_executable.sh"),
			},
			wantExitCode: 127,
			wantErr:      true,
		},
		{
			name: "exit 42",
			args: args{
				cmd: exec.Command("bash", "testdata/exit_42.sh"),
			},
			wantExitCode: 42,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExitCode, err := command.RunCmdAndReturnExitCode(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCmdAndReturnExitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("RunCmdAndReturnExitCode() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
		})
	}
}
