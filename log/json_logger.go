package log

import (
	"fmt"
	"io"
	"os"
)

// JSONLoger ...
type JSONLoger struct {
	writer io.Writer
}

// NewJSONLoger ...
func NewJSONLoger(writer io.Writer) *JSONLoger {
	return &JSONLoger{
		writer: writer,
	}
}

// NewDefaultJSONLoger ...
func NewDefaultJSONLoger() JSONLoger {
	return JSONLoger{
		writer: os.Stdout,
	}
}

// PrintO ...
func (l JSONLoger) PrintO(f Formatable) {
	fmt.Fprintln(l.writer, f.JSON())
}

// Printf ...
func (l JSONLoger) Printf(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	l.PrintO(Message{Content: str})
}
