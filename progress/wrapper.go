package progress

import (
	"io"
	"os"
	"strings"
)

// Wrapper wraps an action with progress indication.
type Wrapper struct {
	ticker *Ticker
}

// newWrapper creates a Wrapper with the given ticker.
func newWrapper(ticker *Ticker) Wrapper {
	return Wrapper{
		ticker: ticker,
	}
}

// NewDefaultWrapper creates a Wrapper with default ticker configuration.
func NewDefaultWrapper(message string) Wrapper {
	return NewDefaultWrapperWithOutput(message, os.Stdout)
}

// NewDefaultWrapperWithOutput creates a Wrapper with default ticker configuration using custom output.
func NewDefaultWrapperWithOutput(message string, output io.Writer) Wrapper {
	// Add "..." suffix if not present
	if !strings.HasSuffix(message, ".") {
		message = message + "..."
	}
	ticker := NewDefaultTicker(message, output)
	return newWrapper(ticker)
}

// WrapAction executes the given action with progress indication.
func (w Wrapper) WrapAction(action func()) {
	w.ticker.Start()
	action()
	w.ticker.Stop()
}
