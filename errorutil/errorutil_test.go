package errorutil_test

import (
	"os/exec"
	"testing"

	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/stretchr/testify/require"
)

func TestIsExitStatusErrorStr(t *testing.T) {
	// --- Should match ---
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 1"))
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 0"))
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 2"))
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 11"))
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 111"))
	require.Equal(t, true, errorutil.IsExitStatusErrorStr("exit status 999"))

	// --- Should not match ---
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("xit status 1"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("status 1"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status "))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 2112"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 21121"))

	// prefixed
	require.Equal(t, false, errorutil.IsExitStatusErrorStr(".exit status 1"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr(" exit status 1"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("error: exit status 1"))
	// postfixed
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 1."))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 1 "))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 1 - something else"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("exit status 1 2"))

	// other
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("-exit status 211-"))
	require.Equal(t, false, errorutil.IsExitStatusErrorStr("something else: exit status 1"))
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
			want: -1,
		},
		{
			name: "invalid command",
			cmd:  exec.Command(""),
			want: -1,
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

			// Act
			got, err := errorutil.CmdExitCodeFromError(err)

			// Assert
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIsExitStatusError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "exit 42",
			err:  exec.Command("bash", "testdata/exit_42.sh").Run(),
			want: true,
		},
		{
			name: "invalid command",
			err:  exec.Command("").Run(),
			want: false,
		},
		{
			name: "nil test",
			err:  nil,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errorutil.IsExitStatusError(tt.err)
			require.Equal(t, tt.want, got)
		})
	}
}
