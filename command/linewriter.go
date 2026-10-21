package command

import (
	"bytes"
	"io"
)

// maxLineLength is the longest line a matcher is offered. A longer one is dropped rather than cut
// short: a reason has to be a short summary, so a line this long is not one, and half of it matching
// would put the whole line in the error that is stored and shown. The stream still forwards it whole.
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

	pending  []byte
	skipping bool
	closed   bool
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

	if len(w.pending) > 0 || w.skipping {
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
	if w.skipping {
		return
	}

	if len(w.pending)+len(b) > maxLineLength {
		w.pending = w.pending[:0]
		w.skipping = true

		return
	}

	w.pending = append(w.pending, b...)
}

// emitLine hands over a line that ended in a newline, without the carriage return of a CRLF: the two
// bytes are one terminator, the way bufio.ScanLines reads them.
func (w *lineWriter) emitLine() {
	w.emit(bytes.TrimSuffix(w.pending, carriageReturn))
}

// emit hands the line over. A dropped one is handed over empty rather than not at all, because a
// matcher reading a multi-line block still has to see that the block ended here.
func (w *lineWriter) emit(line []byte) {
	if w.skipping {
		line = nil
	}

	w.collector.consumeLine(line)
	w.pending = w.pending[:0]
	w.skipping = false
}
