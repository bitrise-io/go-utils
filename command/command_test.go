package command

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	"github.com/stretchr/testify/assert"
)

func TestRunErrors(t *testing.T) {
	tests := []struct {
		name    string
		cmd     command
		wantErr string
	}{
		{
			name:    "command without stdout set",
			cmd:     command{cmd: exec.Command("bash", "testdata/exit_with_message.sh")},
			wantErr: `command failed with exit status 1 (bash "testdata/exit_with_message.sh"): check the command's output for details`,
		},
		{
			name: "command with stdout set",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				var out bytes.Buffer
				c.Stdout = &out
				return command{cmd: c}
			}(),
			wantErr: `command failed with exit status 1 (bash "testdata/exit_with_message.sh"): check the command's output for details`,
		},
		{
			name: "command with error finder",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				errorFinder := func(out string) []string {
					var errors []string
					for _, line := range strings.Split(out, "\n") {
						if strings.Contains(line, "Error:") {
							errors = append(errors, line)
						}
					}
					return errors
				}

				return command{
					cmd:            c,
					errorCollector: &errorCollector{errorFinder: errorFinder},
				}
			}(),
			wantErr: `command failed with exit status 1 (bash "testdata/exit_with_message.sh"): Error: first error
Error: second error
Error: third error
Error: fourth error`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Run()
			var gotErrMsg string
			if err != nil {
				gotErrMsg = err.Error()
			}
			if gotErrMsg != tt.wantErr {
				t.Errorf("command.Run() error = \n%v\n, wantErr \n%v\n", gotErrMsg, tt.wantErr)
				return
			}

			gotFormattedMsg := errorutil.FormattedError(err)
			if gotFormattedMsg != tt.wantErr {
				t.Errorf("FormattedError() error = \n%v\n, wantErr \n%v\n", gotFormattedMsg, tt.wantErr)
			}
		})
	}
}

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
			if tt.wantErr && tt.wantExitCode > 0 {
				var exitErr *exec.ExitError

				if ok := errors.As(err, &exitErr); !ok {
					t.Errorf("command.RunAndReturnExitCode() did nor return ExitError type: %s", err)
					return
				}

				if exitErr.ExitCode() != tt.wantExitCode {
					t.Errorf("command.RunAndReturnExitCode() exit code = %v, want %v", exitErr.ExitCode(), tt.wantExitCode)
				}
			}
			if gotExitCode != tt.wantExitCode {
				t.Errorf("command.RunAndReturnExitCode() = %v, want %v", gotExitCode, tt.wantExitCode)
			}
		})
	}
}

func TestRunAndReturnTrimmedOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     command
		wantErr string
	}{
		{
			name: "command without error finder",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				return command{
					cmd: c,
				}
			}(),
			wantErr: "command failed with exit status 1 (bash \"testdata/exit_with_message.sh\"): check the command's output for details",
		},
		{
			name: "command with error finder",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				errorFinder := func(out string) []string {
					var errors []string
					for _, line := range strings.Split(out, "\n") {
						if strings.Contains(line, "Error:") {
							errors = append(errors, line)
						}
					}
					return errors
				}

				return command{
					cmd:            c,
					errorCollector: &errorCollector{errorFinder: errorFinder},
				}
			}(),
			wantErr: `command failed with exit status 1 (bash "testdata/exit_with_message.sh"): Error: first error
Error: second error`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.cmd.RunAndReturnTrimmedOutput()
			var gotErrMsg string
			if err != nil {
				gotErrMsg = err.Error()
			}
			if gotErrMsg != tt.wantErr {
				t.Errorf("command.Run() error = %v, wantErr %v", gotErrMsg, tt.wantErr)
				return
			}
		})
	}
}

func TestRunAndReturnTrimmedCombinedOutput(t *testing.T) {
	tests := []struct {
		name    string
		cmd     command
		wantErr string
	}{
		{
			name: "command without error finder",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				return command{
					cmd: c,
				}
			}(),
			wantErr: "command failed with exit status 1 (bash \"testdata/exit_with_message.sh\"): check the command's output for details",
		},
		{
			name: "command with error finder",
			cmd: func() command {
				c := exec.Command("bash", "testdata/exit_with_message.sh")
				errorFinder := func(out string) []string {
					var errors []string
					for _, line := range strings.Split(out, "\n") {
						if strings.Contains(line, "Error:") {
							errors = append(errors, line)
						}
					}
					return errors
				}

				return command{
					cmd:            c,
					errorCollector: &errorCollector{errorFinder: errorFinder},
				}
			}(),
			wantErr: `command failed with exit status 1 (bash "testdata/exit_with_message.sh"): Error: first error
Error: second error
Error: third error
Error: fourth error`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.cmd.RunAndReturnTrimmedCombinedOutput()
			var gotErrMsg string
			if err != nil {
				gotErrMsg = err.Error()
			}
			if gotErrMsg != tt.wantErr {
				t.Errorf("command.Run() error = %v, wantErr %v", gotErrMsg, tt.wantErr)
				return
			}
		})
	}
}

func TestSpecialCharactersAreNotEscaped(t *testing.T) {
	programName := "test"
	argument := `-----BEGIN PRIVATE KEY-----\nThis\nis\na\nprivate-key\n-----END PRIVATE KEY-----`

	got := NewFactory(env.NewRepository()).Create(programName, []string{argument}, nil).PrintableCommandArgs()
	expected := fmt.Sprintf("%s \"%s\"", programName, argument)

	assert.Equal(t, expected, got)
}
