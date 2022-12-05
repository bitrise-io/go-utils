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
