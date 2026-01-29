package progress

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal" //nolint:staticcheck // Keep using this for Go 1.21 compatibility
)

// Wrapper wraps an action with progress indication.
type Wrapper struct {
	spinner         Spinner
	interactiveMode bool
}

// NewWrapper creates a Wrapper with the given spinner and interactive mode setting.
func NewWrapper(spinner Spinner, interactiveMode bool) Wrapper {
	return Wrapper{
		spinner:         spinner,
		interactiveMode: interactiveMode,
	}
}

// NewDefaultWrapper creates a Wrapper with default spinner configuration.
func NewDefaultWrapper(message string) Wrapper {
	spinner := NewDefaultSpinner(message)
	interactiveMode := OutputDeviceIsTerminal()
	return NewWrapper(spinner, interactiveMode)
}

// NewDefaultWrapperWithOutput creates a Wrapper with default spinner configuration.
func NewDefaultWrapperWithOutput(message string, output io.Writer) Wrapper {
	spinner := NewDefaultSpinnerWithOutput(message, output)
	interactiveMode := OutputDeviceIsTerminal()
	return NewWrapper(spinner, interactiveMode)
}

// WrapAction executes the given action with progress indication.
func (w Wrapper) WrapAction(action func()) {
	if w.interactiveMode {
		w.spinner.Start()
		action()
		w.spinner.Stop()
	} else {
		message := w.spinner.message
		if !strings.HasSuffix(message, ".") {
			message = message + "..."
		}
		if _, err := fmt.Fprintln(w.spinner.writer, message); err != nil {
			fmt.Printf("failed to print message: %s, error: %s", message, err)
		}
		action()
	}
}

// OutputDeviceIsTerminal returns true if stdout is connected to a terminal.
func OutputDeviceIsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
