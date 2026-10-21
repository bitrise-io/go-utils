package command

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// maxLineLength bounds how much of a line the matcher is given. The stream still forwards a longer
// line whole; only what the matcher reads is cut, so that one unbroken megabyte of output cannot be
// held in memory on its way to it.
const maxLineLength = 64 * 1024

var (
	newline        = []byte{'\n'}
	carriageReturn = []byte{'\r'}
)

// lineWriter forwards everything written to it and assembles the lines behind it, so that the matcher
// is given complete lines however the output was chunked. An escape sequence split across two writes
// is reassembled with the line it belongs to.
//
// Each stream gets its own lineWriter, because a partial line has to be held per stream.
type lineWriter struct {
	dest      io.Writer
	collector *matcherCollector

	pending []byte
	closed  bool
}

// Write forwards p untouched and feeds the lines it completes to the matcher. Output reaches the
// caller's writer as it arrives, so a line a command leaves unfinished is on screen while it runs.
func (w *lineWriter) Write(p []byte) (int, error) {
	n, err := w.dest.Write(p)
	w.consume(p[:n])

	return n, err
}

// Close hands the matcher a last line that has no newline, which is how a command that ends mid-line
// says it has nothing more to write. It does not close the stream's own writer.
func (w *lineWriter) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	if len(w.pending) > 0 {
		w.emit(w.pending)
	}

	return nil
}

func (w *lineWriter) consume(p []byte) {
	for len(p) > 0 {
		line, rest, found := bytes.Cut(p, newline)
		w.hold(line)
		if !found {
			return
		}

		w.emitLine()
		p = rest
	}
}

func (w *lineWriter) hold(b []byte) {
	room := maxLineLength - len(w.pending)
	if len(b) <= room {
		w.pending = append(w.pending, b...)
		return
	}

	w.pending = append(w.pending, trimPartialRune(b[:room])...)
}

// emitLine hands over a line that ended in a newline, without the carriage return of a CRLF: the two
// bytes are one terminator, the way bufio.ScanLines reads them.
func (w *lineWriter) emitLine() {
	w.emit(bytes.TrimSuffix(w.pending, carriageReturn))
}

func (w *lineWriter) emit(line []byte) {
	w.collector.consumeLine(line)
	w.pending = w.pending[:0]
}

// trimPartialRune drops a multi-byte rune that the cut would leave half written, so that the matcher
// is never handed an invalid one.
func trimPartialRune(b []byte) []byte {
	for i := len(b); i > 0 && i > len(b)-utf8.UTFMax; i-- {
		if r, size := utf8.DecodeLastRune(b[:i]); r != utf8.RuneError || size > 1 {
			return b[:i]
		}
	}

	return b
}
