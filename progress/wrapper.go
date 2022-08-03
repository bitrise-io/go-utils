package progress

import (
	"fmt"
	"strings"
)

// Wrapper ...
type Wrapper struct {
	spinner         Spinner
	action          func()
	interactiveMode bool
}

// NewWrapper ...
func NewWrapper(spinner Spinner, interactiveMode bool) Wrapper {
	return Wrapper{
		spinner:         spinner,
		interactiveMode: interactiveMode,
	}
}

// WrapAction ...
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
