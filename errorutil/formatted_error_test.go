package errorutil

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

func TestFormattedError(t *testing.T) {
	tests := []struct {
		name               string
		error              func() error
		wantFormattedError string
	}{
		{
			name: "Single error",
			error: func() error {
				return errors.New("this is a failed error")
			},
			wantFormattedError: "this is a failed error",
		},
		{
			name: "Single multiline error",
			error: func() error {
				return errors.New("this is a failed error\nAnother line\nanother line")
			},
			wantFormattedError: "this is a failed error\nAnother line\nanother line",
		},
		{
			name: "Multiple wrapped errors",
			error: func() error {
				err := errors.New("the magic has failed")
				err = fmt.Errorf("second layer also failed: %w", err)
				err = fmt.Errorf("third layer also failed: %w", err)
				err = fmt.Errorf("fourth layer also failed: %w", err)
				return err
			}, wantFormattedError: `fourth layer also failed:
  third layer also failed:
    second layer also failed:
      the magic has failed`,
		},
		{
			name: "Multiple non-wrapped errors",
			error: func() error {
				err := errors.New("the magic has failed")
				err = fmt.Errorf("second layer also failed: %s", err)
				err = fmt.Errorf("third layer also failed: %s", err)
				err = fmt.Errorf("fourth layer also failed: %s", err)
				return err
			}, wantFormattedError: "fourth layer also failed: third layer also failed: second layer also failed: the magic has failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormattedError(tt.error())

			if formatted != tt.wantFormattedError {
				t.Errorf("got formatted error = %s, want %s", formatted, tt.wantFormattedError)
			}
		})
	}
}

func TestFormattedErrorWithCommand(t *testing.T) {
	commandFactory := command.NewFactory(env.NewRepository())

	tests := []struct {
		name    string
		cmdFn   func() error
		wantErr string
		wantMsg string
	}{
		{
			name: "command without stdout set",
			cmdFn: func() error {
				cmd := commandFactory.Create("bash", []string{"../command/testdata/exit_with_message.sh"}, nil)
				return cmd.Run()
			},
			wantErr: `command failed with exit status 1 (bash "../command/testdata/exit_with_message.sh"): check the command's output for details`,
			wantMsg: `command failed with exit status 1 (bash "../command/testdata/exit_with_message.sh"):
  check the command's output for details`,
		},
		{
			name: "command with error finder",
			cmdFn: func() error {
				errorFinder := func(out string) []string {
					var errors []string
					for _, line := range strings.Split(out, "\n") {
						if strings.Contains(line, "Error:") {
							errors = append(errors, line)
						}
					}
					return errors
				}

				cmd := commandFactory.Create("bash", []string{"../command/testdata/exit_with_message.sh"}, &command.Opts{
					ErrorFinder: errorFinder,
				})

				err := cmd.Run()
				return err
			},
			wantErr: `command failed with exit status 1 (bash "../command/testdata/exit_with_message.sh"): Error: first error
Error: second error
Error: third error
Error: fourth error`,
			wantMsg: `command failed with exit status 1 (bash "../command/testdata/exit_with_message.sh"):
  Error: first error
  Error: second error
  Error: third error
  Error: fourth error`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmdFn()

			var gotErrMsg string
			if err != nil {
				gotErrMsg = err.Error()
			}
			if gotErrMsg != tt.wantErr {
				t.Errorf("command.Run() error = \n%v\n, wantErr \n%v\n", gotErrMsg, tt.wantErr)
				return
			}

			gotFormattedMsg := FormattedError(err)
			require.Equal(t, tt.wantMsg, gotFormattedMsg, "FormattedError() error = \n%v\n, wantErr \n%v\n", gotFormattedMsg, tt.wantErr)
			if gotFormattedMsg != tt.wantMsg {
				t.Errorf("FormattedError() error = \n%v\n, wantErr \n%v\n", gotFormattedMsg, tt.wantErr)
			}
		})
	}
}
