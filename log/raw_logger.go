package log

import (
	"fmt"
	"io"
	"os"
)

// RawLogger ...
type RawLogger struct {
	writer io.Writer
}

// NewRawLogger ...
func NewRawLogger(writer io.Writer) *RawLogger {
	return &RawLogger{
		writer: writer,
	}
}

// NewDefaultRawLogger ...
func NewDefaultRawLogger() RawLogger {
	return RawLogger{
		writer: os.Stdout,
	}
}

// Printd ...
func (l RawLogger) Printd(f Formatable) {
	fmt.Fprintln(l.writer, f.String())
}

// Printf ...
func (l RawLogger) Printf(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	l.Printd(Message{Content: str})
}
