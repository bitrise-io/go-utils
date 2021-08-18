package command

import (
	"os/exec"
	"testing"
)

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
			command := NewCommand(tt.args.cmd)
			gotExitCode, err := command.RunAndReturnExitCode()
			if (err != nil) != tt.wantErr {
				t.Errorf("runCmdAndReturnExitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("runCmdAndReturnExitCode() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
		})
	}
}
