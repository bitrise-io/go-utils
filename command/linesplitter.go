package command

import "bytes"

// maxLineLength is the longest line a matcher is offered. A longer one is dropped rather than cut
// short, so that a matcher only ever sees complete lines: a line this long is not a failure reason,
// and a cut one that matched would fill the reason with noise. The stream still forwards it whole.
const maxLineLength = 64 * 1024

var (
	newline        = []byte{'\n'}
	carriageReturn = []byte{'\r'}
)

// lineSplitter is the observer that assembles lines, so that the matcher is given complete lines
// however the output was chunked. An escape sequence split across two writes is reassembled with
// the line it belongs to. It forwards nothing: the tap in front of it has already done that.
type lineSplitter struct {
	matcher ErrorMatcher

	pending  []byte
	skipping bool
	closed   bool
}

func (s *lineSplitter) observe(p []byte) {
	for len(p) > 0 {
		line, rest, found := bytes.Cut(p, newline)
		s.hold(line)
		if !found {
			return
		}

		s.emitLine()
		p = rest
	}
}

// close hands the matcher a last line that has no newline, which is how a command that ends
// mid-line says it has nothing more to write.
func (s *lineSplitter) close() {
	if s.closed {
		return
	}
	s.closed = true

	if len(s.pending) > 0 || s.skipping {
		s.emit(s.pending)
	}
}

func (s *lineSplitter) hold(b []byte) {
	if s.skipping {
		return
	}

	if len(s.pending)+len(b) > maxLineLength {
		s.pending = s.pending[:0]
		s.skipping = true

		return
	}

	s.pending = append(s.pending, b...)
}

// emitLine hands over a line that ended in a newline, without the carriage return of a CRLF: the two
// bytes are one terminator, the way bufio.ScanLines reads them.
func (s *lineSplitter) emitLine() {
	s.emit(bytes.TrimSuffix(s.pending, carriageReturn))
}

// emit hands the line over. A dropped one is handed over empty rather than not at all, because a
// matcher reading a multi-line block still has to see that the block ended here.
func (s *lineSplitter) emit(line []byte) {
	if s.skipping {
		line = nil
	}

	s.matcher.ConsumeLine(string(line))
	s.pending = s.pending[:0]
	s.skipping = false
}
