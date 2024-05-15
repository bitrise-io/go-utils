package errorutil

import (
	"errors"
	"fmt"
	"testing"
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
		{
			name: "Multiple wrapped errors in the same level",
			error: func() error {
				err := errors.New("the internal debug info")
				err2 := fmt.Errorf("the description")
				err = fmt.Errorf("third layer also failed: %w %w", err2, err)
				err = fmt.Errorf("fourth layer also failed: %w", err)
				return err
			}, wantFormattedError: `fourth layer also failed:
  third layer also failed:
    the description
    the internal debug info`,
		},
		{
			name: "Multiple wrapped errors in the same level, debug info hidden from stack trace",
			error: func() error {
				err := NewInternalDebugError(errors.New("the internal debug info"))
				err2 := fmt.Errorf("the description")
				err = fmt.Errorf("third layer also failed: %w %w", err, err2)
				err = fmt.Errorf("fourth layer also failed: %w", err)
				return err
			}, wantFormattedError: `fourth layer also failed:
  third layer also failed:
    the description`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormattedError(tt.error())

			// require.Equal(t, tt.wantFormattedError, formatted)
			if formatted != tt.wantFormattedError {
				t.Errorf("got formatted error = \n%s\n, want \n%s", formatted, tt.wantFormattedError)
			}
		})
	}
}
