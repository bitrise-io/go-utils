package errorutil

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsExitStatusErrorStr(t *testing.T) {
	// --- Should match ---
	require.Equal(t, true, IsExitStatusErrorStr("exit status 1"))
	require.Equal(t, true, IsExitStatusErrorStr("exit status 0"))
	require.Equal(t, true, IsExitStatusErrorStr("exit status 2"))
	require.Equal(t, true, IsExitStatusErrorStr("exit status 11"))
	require.Equal(t, true, IsExitStatusErrorStr("exit status 111"))
	require.Equal(t, true, IsExitStatusErrorStr("exit status 999"))

	// --- Should not match ---
	require.Equal(t, false, IsExitStatusErrorStr("xit status 1"))
	require.Equal(t, false, IsExitStatusErrorStr("status 1"))
	require.Equal(t, false, IsExitStatusErrorStr("exit status "))
	require.Equal(t, false, IsExitStatusErrorStr("exit status"))
	require.Equal(t, false, IsExitStatusErrorStr("exit status 2112"))
	require.Equal(t, false, IsExitStatusErrorStr("exit status 21121"))

	// prefixed
	require.Equal(t, false, IsExitStatusErrorStr(".exit status 1"))
	require.Equal(t, false, IsExitStatusErrorStr(" exit status 1"))
	require.Equal(t, false, IsExitStatusErrorStr("error: exit status 1"))
	// postfixed
	require.Equal(t, false, IsExitStatusErrorStr("exit status 1."))
	require.Equal(t, false, IsExitStatusErrorStr("exit status 1 "))
	require.Equal(t, false, IsExitStatusErrorStr("exit status 1 - something else"))
	require.Equal(t, false, IsExitStatusErrorStr("exit status 1 2"))

	// other
	require.Equal(t, false, IsExitStatusErrorStr("-exit status 211-"))
	require.Equal(t, false, IsExitStatusErrorStr("something else: exit status 1"))
}

func TestCmdExitCodeFromError(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		cmd  *exec.Cmd
		want int
	}{
		{
			name: "env command",
			cmd:  exec.Command("env"),
			want: 0,
		},
		{
			name: "not existing executable",
			cmd:  exec.Command("bash", "testdata/not_existing_executable.sh"),
			want: 127,
		},
		{
			name: "exit 42",
			cmd:  exec.Command("bash", "testdata/exit_42.sh"),
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Run()
			want := tt.cmd.ProcessState.ExitCode()
			if tt.want != want {
				t.Errorf("Invalid test expectation: tt.cmd.ProcessState.ExitCode() = %v, want %v", want, tt.want)
			}

			// Act
			got, err := CmdExitCodeFromError(err)

			// Assert
			if got != want {
				t.Errorf("CmdExitCodeFromError() = %v, want %v", got, want)
			}
		})
	}
}

func TestIsExitStatusError(t *testing.T) {
	tests := []struct {
		name string
		cmd  *exec.Cmd
		want bool
	}{
		{
			name: "exit 42",
			cmd:  exec.Command("bash", "testdata/exit_42.sh"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Run()
			want := tt.cmd.ProcessState.Exited()
			if tt.want != want {
				t.Errorf("Invalid test expectation: tt.cmd.ProcessState.Exited() = %v, want %v", want, tt.want)
			}

			if got := IsExitStatusError(err); got != tt.want {
				t.Errorf("IsExitStatusError() = %v, want %v", got, tt.want)
			}
		})
	}
}
