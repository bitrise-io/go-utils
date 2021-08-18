package command

import (
	"github.com/bitrise-io/go-utils/env"
	"testing"
)

func TestRunCmdAndReturnExitCode(t *testing.T) {
	type args struct {
		cmd Command
	}
	factory := NewFactory(env.NewRepository())
	tests := []struct {
		name         string
		args         args
		wantExitCode int
		wantErr      bool
	}{
		{
			name: "invalid command",
			args: args{
				cmd: factory.Create("", nil, nil),
			},
			wantExitCode: -1,
			wantErr:      true,
		},
		{
			name: "env command",
			args: args{
				cmd: factory.Create("env", nil, nil),
			},
			wantExitCode: 0,
			wantErr:      false,
		},
		{
			name: "not existing executable",
			args: args{
				cmd: factory.Create("bash", []string{"testdata/not_existing_executable.sh"}, nil),
			},
			wantExitCode: 127,
			wantErr:      true,
		},
		{
			name: "exit 42",
			args: args{
				cmd: factory.Create("bash", []string{"testdata/exit_42.sh"}, nil),
			},
			wantExitCode: 42,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExitCode, err := tt.args.cmd.RunAndReturnExitCode()
			if (err != nil) != tt.wantErr {
				t.Errorf("command.RunAndReturnExitCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("command.RunAndReturnExitCode() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
		})
	}
}
