package progress

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFmtPrinter(t *testing.T) {
	t.Run("PrintWithoutNewline outputs to stdout", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		printer := NewFmtPrinter()
		printer.PrintWithoutNewline("test")

		_ = w.Close() // err here might break the test, that should be fine though.
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) // err here might break the test, that should be fine though.

		require.Equal(t, "test", buf.String())
	})

	t.Run("Println outputs newline to stdout", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		printer := NewFmtPrinter()
		printer.Println()

		_ = w.Close() // err here might break the test, that should be fine though.
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) // err here might break the test, that should be fine though.

		require.Equal(t, "\n", buf.String())
	})

	t.Run("combined output", func(t *testing.T) {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		printer := NewFmtPrinter()
		printer.PrintWithoutNewline(".")
		printer.PrintWithoutNewline(".")
		printer.PrintWithoutNewline(".")
		printer.Println()

		_ = w.Close() // err here might break the test, that should be fine though.
		os.Stdout = old

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) // err here might break the test, that should be fine though.

		require.Equal(t, "...\n", buf.String())
	})
}

// MockPrinter is a test implementation of Printer that captures output.
type MockPrinter struct {
	Output string
}

// NewMockPrinter creates a new MockPrinter.
func NewMockPrinter() *MockPrinter {
	return &MockPrinter{}
}

// PrintWithoutNewline appends text to the output buffer.
func (m *MockPrinter) PrintWithoutNewline(text string) {
	m.Output += text
}

// Println appends a newline to the output buffer.
func (m *MockPrinter) Println() {
	m.Output += "\n"
}

func TestMockPrinter(t *testing.T) {
	printer := NewMockPrinter()

	printer.PrintWithoutNewline("Hello")
	printer.PrintWithoutNewline(" ")
	printer.PrintWithoutNewline("World")
	printer.Println()

	require.Equal(t, "Hello World\n", printer.Output)
}
