// Package linepipe assembles the lines of a command's output out of the chunks it arrives in and runs
// a set of stages over each complete line, so that a stage carries its rules only and not a buffering
// scheme of its own.
//
// An escape sequence split across two writes is reassembled with the line it belongs to, so a stage
// never sees half of one.
package linepipe

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// MaxLineLength bounds how much of a line the stages are given. The stream still forwards a longer
// line whole; only what the stages read is cut, so that one unbroken megabyte of output cannot be
// held in memory on its way to them.
const MaxLineLength = 64 * 1024

var newline = []byte{'\n'}

// Option configures a Writer.
type Option func(*Writer)

// Writer assembles lines out of the chunks written to it, runs its stages over each complete line,
// and forwards every byte to its destination.
//
// A Writer serves one stream, so give stdout and stderr one each unless the command itself shares a
// single writer between them.
type Writer struct {
	dest io.Writer
	taps []Tap

	pending []byte
	closed  bool
	line    Line
}

// New returns a Writer that forwards to dest and runs opts' stages over each line. A nil dest
// discards the output, so a caller that only wants the stages does not need a writer for it.
func New(dest io.Writer, opts ...Option) *Writer {
	if dest == nil {
		dest = io.Discard
	}

	w := &Writer{dest: dest}
	for _, opt := range opts {
		opt(w)
	}

	return w
}

// WithTap appends t to the stages the Writer runs. Stages run in the order they are given, and that
// order is a correctness rule as soon as a stage rewrites a line: a tap that records text has to come
// after the one that redacts it, or a secret reaches what it records.
func WithTap(t Tap) Option {
	return func(w *Writer) {
		w.taps = append(w.taps, t)
	}
}

// Write forwards p untouched and feeds the lines it completes to the stages. Output reaches the
// destination as it arrives, so a line a command leaves unfinished is still on screen while it runs.
func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.dest.Write(p)
	w.consume(p[:n])

	return n, err
}

// Close hands the stages a last line that has no newline, which is how a command that ends mid-line
// says it has nothing more to write. It does not close the destination.
func (w *Writer) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	if len(w.pending) > 0 {
		w.emit()
	}

	return nil
}

func (w *Writer) consume(p []byte) {
	for len(p) > 0 {
		line, rest, found := bytes.Cut(p, newline)
		w.hold(line)
		if !found {
			return
		}

		w.emit()
		p = rest
	}
}

// hold adds b to the line being assembled, up to MaxLineLength. Past that the line is cut rather than
// buffered: the stages read lines, and nothing a stage looks for is that long.
func (w *Writer) hold(b []byte) {
	room := MaxLineLength - len(w.pending)
	if len(b) <= room {
		w.pending = append(w.pending, b...)
		return
	}

	w.pending = append(w.pending, trimPartialRune(b[:room])...)
}

func (w *Writer) emit() {
	w.line.reset(w.pending)
	for _, tap := range w.taps {
		tap.Observe(&w.line)
	}

	w.pending = w.pending[:0]
}

// trimPartialRune drops a multi-byte rune that the cut would leave half written, so that a stage
// reading the line as text is never handed an invalid one.
func trimPartialRune(b []byte) []byte {
	for i := len(b); i > 0 && i > len(b)-utf8.UTFMax; i-- {
		if r, size := utf8.DecodeLastRune(b[:i]); r != utf8.RuneError || size > 1 {
			return b[:i]
		}
	}

	return b
}
